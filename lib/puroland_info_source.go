package scraper

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

const (
	PurolandInfoUrl    = "http://www.puroland.jp/"
	PurolandInfoApiUrl = "http://www.puroland.jp/api/live/get_information/?page=1&count=20"
)

type information struct {
	Status  string            `json:"status"`
	Count   int               `json:"count"`
	Total   int               `json:"total"`
	Page    int               `json:"page"`
	MaxPage int               `json:"maxpage"`
	Data    []informationItem `json:"data"`
}

type informationItem struct {
	Url             string `json:"url"`
	Title           string `json:"title"`
	PublicDate      string `json:"public_date"`
	ThumbnailMiddle string `json:"thumbnail_m"`
}

type PurolandInfoSource struct {
}

func NewPurolandInfoSource() *PurolandInfoSource {
	return &PurolandInfoSource{}
}

func (s *PurolandInfoSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get(PurolandInfoApiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return s.ScrapeFromReader(res.Body)
}

func (s *PurolandInfoSource) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	jsonData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var info information
	if err := json.Unmarshal(jsonData, &info); err != nil {
		return nil, err
	}

	items := make([]*feeds.Item, 0, info.Count)
	for _, infoItem := range info.Data {
		if infoItem.PublicDate == "" {
			continue
		}
		created, err := time.Parse("2006/01/02", infoItem.PublicDate)
		if err != nil {
			return nil, err
		}

		description := fmt.Sprintf(`<img src="%s" />`, infoItem.ThumbnailMiddle)

		items = append(items, &feeds.Item{
			Title:       html.UnescapeString(infoItem.Title),
			Link:        &feeds.Link{Href: infoItem.Url},
			Id:          infoItem.Url,
			Created:     created,
			Description: description,
		})
	}

	feed := &feeds.Feed{
		Title: "お知らせ | サンリオピューロランド",
		Link:  &feeds.Link{Href: PurolandInfoUrl},
		Items: items,
	}

	return feed, nil
}