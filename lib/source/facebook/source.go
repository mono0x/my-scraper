package facebook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/pkg/errors"
)

const (
	serviceURL = "https://www.facebook.com/"
	baseURL    = "https://graph.facebook.com"
)

type posts struct {
	// https://developers.facebook.com/docs/graph-api/reference/v2.6/post
	Data []struct {
		Id          string `json:"id"`
		CreatedTime string `json:"created_time"`
		From        struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"from"`
		Link    string `json:"link"`
		Message string `json:"message"`
		Picture string `json:"picture"`
	} `json:"data"`
}

type source struct {
	httpClient  *http.Client
	accessToken string
	userID      string
	baseURL     string // for testing
}

var _ scraper.Source = (*source)(nil)

func NewSource(c *http.Client, accessToken string, userID string) *source {
	return &source{
		httpClient:  c,
		accessToken: accessToken,
		userID:      userID,
		baseURL:     baseURL,
	}
}

var (
	photosURLRe     = regexp.MustCompile(`^` + regexp.QuoteMeta(serviceURL) + `[^/]+/photos/`)
	messageReplacer = strings.NewReplacer("\n", "<br />")
)

func (s *source) Scrape() (*feeds.Feed, error) {
	posts, err := s.fetch()
	if err != nil {
		return nil, err
	}
	return s.render(posts)
}

func (s *source) fetch() (*posts, error) {
	values := &url.Values{}
	values.Set("access_token", s.accessToken)
	values.Set("fields", "created_time,from,link,message,picture")

	resp, err := s.httpClient.Get(s.baseURL + "/v2.6/" + s.userID + "/posts?" + values.Encode())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	var posts posts
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, err
	}
	return &posts, nil
}

func (s *source) render(posts *posts) (*feeds.Feed, error) {
	if len(posts.Data) == 0 {
		return nil, errors.New("no posts found")
	}

	items := make([]*feeds.Item, 0, len(posts.Data))
	for _, post := range posts.Data {
		created, err := time.Parse("2006-01-02T15:04:05-0700", post.CreatedTime)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		var title, description string
		if index := strings.Index(post.Message, "\n"); index >= 0 {
			title = post.Message[0:index]
			description = messageReplacer.Replace(post.Message)
		} else {
			title = post.Message
			description = post.Message
		}
		if post.Picture != "" {
			description += fmt.Sprintf(`<br /><img src="%s" />`, post.Picture)
		}

		var link string
		photosURLRe := photosURLRe.Copy()
		if photosURLRe.MatchString(post.Link) {
			if parts := strings.SplitN(post.Id, "_", 2); len(parts) == 2 {
				link = serviceURL + s.userID + "/posts/" + parts[1] + "/"
			} else {
				link = post.Link
			}
		} else {
			link = post.Link
		}

		items = append(items, &feeds.Item{
			Id:          post.Id,
			Author:      &feeds.Author{Name: post.From.Name},
			Title:       title,
			Description: description,
			Created:     created,
			Link:        &feeds.Link{Href: link},
		})
	}

	feed := &feeds.Feed{
		Title: posts.Data[0].From.Name,
		Link:  &feeds.Link{Href: serviceURL + s.userID},
		Items: items,
	}
	return feed, nil
}
