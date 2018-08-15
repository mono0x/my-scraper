package twitter

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/gorilla/feeds"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/pkg/errors"
)

type source struct {
	httpClient *http.Client
	userID     int64
	baseURL    string // for testing
}

var _ scraper.Source = (*source)(nil)

func NewSource(c *http.Client, userID int64) *source {
	return &source{
		httpClient: c,
		userID:     userID,
		baseURL:    anaconda.BaseUrl,
	}
}

func (s *source) Scrape() (*feeds.Feed, error) {
	api := anaconda.NewTwitterApiWithCredentials(
		os.Getenv("TWITTER_OAUTH_TOKEN"),
		os.Getenv("TWITTER_OAUTH_TOKEN_SECRET"),
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"))
	defer api.Close()

	api.HttpClient = s.httpClient
	api.SetBaseUrl(s.baseURL)

	values := url.Values{}
	values.Set("user_id", strconv.FormatInt(s.userID, 10))
	values.Set("count", "100")
	timeline, err := api.GetUserTimeline(values)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.render(timeline)
}

func (s *source) render(timeline []anaconda.Tweet) (*feeds.Feed, error) {
	if len(timeline) == 0 {
		return nil, errors.New("timeline is empty")
	}
	user := timeline[0].User
	userURL := fmt.Sprintf("https://twitter.com/%s", user.ScreenName)
	items := make([]*feeds.Item, 0, len(timeline))
	for _, tweet := range timeline {
		created, err := tweet.CreatedAtTime()
		if err != nil {
			continue
		}
		items = append(items, &feeds.Item{
			Title:   tweet.Text,
			Created: created,
			Link:    &feeds.Link{Href: fmt.Sprintf("%s/status/%s", userURL, tweet.IdStr)},
		})
	}
	return &feeds.Feed{
		Title: fmt.Sprintf("%s (@%s)", user.Name, user.ScreenName),
		Link:  &feeds.Link{Href: userURL},
		Items: items,
	}, nil
}
