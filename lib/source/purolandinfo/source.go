package purolandinfo

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/gorilla/feeds"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/pkg/errors"
)

const (
	baseURL                 = "http://www.puroland.jp"
	purolandInfoAPIEndpoint = "/api/live/get_information/?page=1&count=20"
)

type source struct {
	httpClient *http.Client
	baseURL    string // for testing
}

var _ scraper.Source = (*source)(nil)

func NewSource(c *http.Client) *source {
	return &source{
		httpClient: c,
		baseURL:    baseURL,
	}
}

func (s *source) Scrape() (*feeds.Feed, error) {
	res, err := s.httpClient.Get(s.baseURL + purolandInfoAPIEndpoint)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	return s.scrapeFromReader(res.Body)
}

func (s *source) scrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	var info struct {
		Status  string `json:"status"`
		Count   int    `json:"count"`
		Total   int    `json:"total"`
		Page    int    `json:"page"`
		MaxPage int    `json:"maxpage"`
		Data    []struct {
			URL             string `json:"url"`
			Title           string `json:"title"`
			PublicDate      string `json:"public_date"`
			ThumbnailMiddle string `json:"thumbnail_m"`
		} `json:"data"`
	}

	if err := json.NewDecoder(reader).Decode(&info); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*feeds.Item, 0, info.Count)
	for _, infoItem := range info.Data {
		if infoItem.PublicDate == "" {
			continue
		}

		description := fmt.Sprintf(`%s<br /><img src="%s" />`, infoItem.PublicDate, infoItem.ThumbnailMiddle)

		items = append(items, &feeds.Item{
			Title:       html.UnescapeString(infoItem.Title),
			Link:        &feeds.Link{Href: infoItem.URL},
			Id:          infoItem.URL,
			Description: description,
		})
	}

	feed := &feeds.Feed{
		Title: "お知らせ | サンリオピューロランド",
		Link:  &feeds.Link{Href: baseURL + "/"},
		Items: items,
	}

	return feed, nil
}
