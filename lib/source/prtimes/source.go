package prtimes

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/pkg/errors"
)

const (
	baseURL  = "https://prtimes.jp"
	endpoint = `/main/action.php?run=html&page=searchkey&search_word=%E3%82%B5%E3%83%B3%E3%83%AA%E3%82%AA&search_pattern=1`
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

	return s.scrapeFromReader(res.Body)
}

func (s *source) scrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.scrapeFromDocument(doc)
}

func (s *source) scrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title: "PR TIMES (Sanrio)",
		Link:  &feeds.Link{Href: baseURL + endpoint},
	}

	absBaseURL, _ := url.Parse(baseURL + endpoint)
	var items []*feeds.Item
	doc.Find("a.link-title-item-ordinary").Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		link, ok := s.Attr("href")
		if !ok {
			return
		}
		refURL, err := url.Parse(link)
		if err != nil {
			return
		}
		absURL := absBaseURL.ResolveReference(refURL)
		link = absURL.String()
		dt, ok := s.Parent().Next().Attr("datetime")
		if !ok {
			return
		}
		t, err := time.Parse("2006-01-02T15:04:05-0700", dt)
		if err != nil {
			return
		}
		author := strings.TrimSpace(s.Parent().Next().Next().Text())
		items = append(items, &feeds.Item{
			Title:   title,
			Link:    &feeds.Link{Href: link},
			Author:  &feeds.Author{Name: author},
			Created: t,
		})
	})
	feed.Items = items

	return feed, nil
}
