package seibuenevent

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	scraper "github.com/mono0x/my-scraper/lib"
)

const (
	baseURL  = "https://www.seibu-leisure.co.jp"
	endpoint = "/event/12410/index.html"
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
	res, err := s.httpClient.Get(s.baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return s.scrapeFromDocument(doc)
}

func (s *source) scrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	var items []*feeds.Item
	var (
		title string
	)
	doc.Find(".elem-section > div > div > div > div > div").Each(func(_ int, s *goquery.Selection) {
		switch {
		case s.HasClass("elem-heading-lv3"):
			title = s.Find("h3").Text()
		case s.HasClass("elem-pic-block"):
			anchor := s.Find("ul.txt-list li p a")
			if anchor.Length() == 0 {
				return
			}
			href, ok := anchor.First().Attr("href")
			if !ok {
				return
			}

			items = append(items, &feeds.Item{
				Title: title,
				Link:  &feeds.Link{Href: baseURL + href},
			})
		}
	})

	feed := &feeds.Feed{
		Title: "西武園ゆうえんち メルヘンタウン",
		Link:  &feeds.Link{Href: baseURL + endpoint},
		Items: items,
	}

	return feed, nil
}
