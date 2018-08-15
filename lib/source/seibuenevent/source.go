package seibuenevent

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/pkg/errors"
)

const (
	baseURL  = "http://www.seibu-leisure.co.jp"
	endpoint = "/event/index.html?category=e1"
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
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
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
			paragraph := s.Find(".elem-paragraph p")
			if paragraph.Length() == 0 {
				return
			}

			description, err := paragraph.Html()
			if err != nil {
				return
			}

			sha := sha256.New()
			fmt.Fprint(sha, title)

			items = append(items, &feeds.Item{
				Title:       title,
				Description: description,
				Link:        &feeds.Link{Href: baseURL + endpoint},
				Id:          fmt.Sprintf("%x", sha.Sum(nil)),
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
