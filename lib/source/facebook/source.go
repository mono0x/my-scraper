package facebook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
)

const (
	ServiceURL  = "https://www.facebook.com/"
	APIEndpoint = "https://graph.facebook.com/v2.6/"
)

type Posts struct {
	Data []*Post `json:"data"`
}

// https://developers.facebook.com/docs/graph-api/reference/v2.6/post
type Post struct {
	Id          string   `json:"id"`
	CreatedTime string   `json:"created_time"`
	From        *Profile `json:"from"`
	Link        string   `json:"link"`
	Message     string   `json:"message"`
	Picture     string   `json:"picture"`
}

type Profile struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Source struct {
	userId string
}

func NewSource(userId string) *Source {
	return &Source{
		userId: userId,
	}
}

var (
	photosURLRe     = regexp.MustCompile(`^` + regexp.QuoteMeta(ServiceURL) + `[^/]+/photos/`)
	messageReplacer = strings.NewReplacer("\n", "<br />")
)

func (s *Source) Scrape() (*feeds.Feed, error) {
	posts, err := s.Fetch()
	if err != nil {
		return nil, err
	}
	return s.Render(posts)
}

func (s *Source) Fetch() (*Posts, error) {
	values := &url.Values{}
	values.Set("access_token", os.Getenv("FACEBOOK_ACCESS_TOKEN"))
	values.Set("fields", "created_time,from,link,message,picture")

	resp, err := http.Get(APIEndpoint + s.userId + "/posts?" + values.Encode())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	var posts Posts
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, err
	}
	return &posts, nil
}

func (s *Source) Render(posts *Posts) (*feeds.Feed, error) {
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
				link = ServiceURL + s.userId + "/posts/" + parts[1] + "/"
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
		Link:  &feeds.Link{Href: ServiceURL + s.userId},
		Items: items,
	}
	return feed, nil
}
