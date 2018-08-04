package instagram

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
)

type InstagramSource struct {
	userId string
}

func NewSource(userId string) *InstagramSource {
	return &InstagramSource{
		userId: userId,
	}
}

func (s *InstagramSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get("https://www.instagram.com/" + s.userId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	return s.ScrapeFromReader(res.Body)
}

var sharedDataRe = regexp.MustCompile(`window\._sharedData\s*=\s*({.+})[\s\n]*[;<]`)

type instagramData struct {
	EntryData struct {
		ProfilePage []struct {
			User struct {
				UserName  string `json:"username"`
				Id        string `json:"id"`
				Biography string `json:"biography"`
				FullName  string `json:"full_name"`
				Media     struct {
					Nodes []struct {
						Code        string `json:"code"`
						Date        int64  `json:"date"`
						Deimensions struct {
							Width  int `json:"width"`
							Height int `json:"height"`
						} `json:"dimensions"`
						Caption      string `json:"caption"`
						ThumbnailSrc string `json:"thumbnail_src"`
						IsVideo      bool   `json:"is_video"`
						Id           string `json:"id"`
						DisplaySrc   string `json:"display_src"`
					} `json:"nodes"`
				} `json:"media"`
			} `json:"user"`
		}
	} `json:"entry_data"`
}

var emojiRe = regexp.MustCompile(`[^\x{0000}-\x{ffff}]+`)

func (s *InstagramSource) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	src, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sharedDataRe := sharedDataRe.Copy()
	submatches := sharedDataRe.FindSubmatch(src)
	if len(submatches) == 0 {
		return nil, errors.New("data not found")
	}

	var data instagramData
	if err := json.Unmarshal(submatches[1], &data); err != nil {
		return nil, errors.WithStack(err)
	}

	if len(data.EntryData.ProfilePage) == 0 {
		return nil, errors.New("ProfilePage item not found")
	}

	user := data.EntryData.ProfilePage[0].User

	items := make([]*feeds.Item, 0, len(user.Media.Nodes))
	for _, node := range user.Media.Nodes {
		caption := emojiRe.ReplaceAllString(node.Caption, "")
		lines := strings.Split(caption, "\n")
		if len(lines) == 0 {
			continue
		}

		title := lines[0]

		escapedLines := make([]string, 0, len(lines))
		for _, line := range lines {
			escapedLines = append(escapedLines, html.EscapeString(line))
		}
		items = append(items, &feeds.Item{
			Title:       title,
			Created:     time.Unix(node.Date, 0).In(time.UTC),
			Link:        &feeds.Link{Href: fmt.Sprintf("http://www.instagram.com/p/%s/", node.Code)},
			Description: fmt.Sprintf("%s<br /><img src=\"%s\" />", strings.Join(escapedLines, "<br />"), node.DisplaySrc),
		})
	}

	return &feeds.Feed{
		Title:       user.FullName,
		Link:        &feeds.Link{Href: fmt.Sprintf("https://www.instagram.com/%s/", user.UserName)},
		Description: user.Biography,
		Items:       items,
	}, nil
}
