package sanrionewsrelease

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/pkg/errors"
)

const (
	baseURL  = "https://www.sanrio.co.jp"
	endpoint = "/corporate/release/"
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

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	doc.Find(".news_release_list dl").Each(func(_ int, s *goquery.Selection) {
		date, err := time.ParseInLocation("2006/1/2", s.Find("dt").First().Text(), loc)
		if err != nil {
			return
		}

		a := s.Find("dd a").First()
		title := strings.TrimSpace(a.Text())
		href, ok := a.Attr("href")
		if !ok {
			return
		}

		if pdf := a.Find("img[alt=PDF]").Length() > 0; pdf {
			title = title + " (PDF)"
		}

		items = append(items, &feeds.Item{
			Title:   title,
			Link:    &feeds.Link{Href: baseURL + href},
			Created: date,
		})
	})

	feed := &feeds.Feed{
		Title: "ニュースリリース | 会社情報 | サンリオ",
		Link:  &feeds.Link{Href: baseURL + endpoint},
		Items: items,
	}
	return feed, nil
}
