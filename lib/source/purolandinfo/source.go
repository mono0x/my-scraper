package purolandinfo

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
)

const (
	purolandInfoURL         = "http://www.puroland.jp/"
	purolandInfoAPIEndpoint = "http://www.puroland.jp/api/live/get_information/?page=1&count=20"
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
	URL             string `json:"url"`
	Title           string `json:"title"`
	PublicDate      string `json:"public_date"`
	ThumbnailMiddle string `json:"thumbnail_m"`
}

type PurolandInfoSource struct {
}

func NewSource() *PurolandInfoSource {
	return &PurolandInfoSource{}
}

func (s *PurolandInfoSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get(purolandInfoAPIEndpoint)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	return s.ScrapeFromReader(res.Body)
}

func (s *PurolandInfoSource) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	jsonData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var info information
	if err := json.Unmarshal(jsonData, &info); err != nil {
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
		Link:  &feeds.Link{Href: purolandInfoURL},
		Items: items,
	}

	return feed, nil
}
