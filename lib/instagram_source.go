package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/feeds"
)

type InstagramSource struct {
	userId string
}

func NewInstagramSource(userId string) *InstagramSource {
	return &InstagramSource{
		userId: userId,
	}
}

func (s *InstagramSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get("https://www.instagram.com/" + s.userId)
	if err != nil {
		return nil, err
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

func (s *InstagramSource) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	src, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	sharedDataRe := sharedDataRe.Copy()
	submatches := sharedDataRe.FindSubmatch(src)
	if len(submatches) == 0 {
		return nil, errors.New("data not found")
	}

	var data instagramData
	if err := json.Unmarshal(submatches[1], &data); err != nil {
		return nil, err
	}

	if len(data.EntryData.ProfilePage) == 0 {
		return nil, errors.New("ProfilePage item not found")
	}

	user := data.EntryData.ProfilePage[0].User

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}

	items := make([]*feeds.Item, 0, len(user.Media.Nodes))
	for _, node := range user.Media.Nodes {
		items = append(items, &feeds.Item{
			Title:       node.Caption,
			Created:     time.Unix(node.Date, 0).In(loc),
			Link:        &feeds.Link{Href: fmt.Sprintf("http://www.instagram.com/p/%s/", node.Code)},
			Description: fmt.Sprintf("%s<br /><img src=\"%s\" />", html.EscapeString(node.Caption), node.DisplaySrc),
		})
	}

	return &feeds.Feed{
		Title:       user.FullName,
		Link:        &feeds.Link{Href: fmt.Sprintf("https://www.instagram.com/%s/", user.UserName)},
		Description: user.Biography,
		Items:       items,
	}, nil
}
