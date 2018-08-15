package valuepress

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/pkg/errors"
)

const (
	baseURL  = "https://www.value-press.com"
	endpoint = `/search?q=%E3%82%B5%E3%83%B3%E3%83%AA%E3%82%AA`
)

type source struct {
	httpClient *http.Client
	baseURL    string
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
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var items []*feeds.Item
	doc.Find(".pressrelease_article").Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2.mb10").First().Text())
		link, ok := s.Find("a").First().Attr("href")
		if !ok {
			return
		}
		description := strings.TrimSpace(s.Find("p").Text())
		t, err := time.ParseInLocation("!2006年1月2日 15時", s.Find(".release_tag li").First().Text(), loc)
		if err != nil {
			return
		}
		author := strings.TrimSpace(s.Find("h3.mt05").First().Text())
		items = append(items, &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: link},
			Description: description,
			Author:      &feeds.Author{Name: author},
			Id:          link,
			Created:     t,
		})
	})

	feed := &feeds.Feed{
		Title: "ValuePress! (Sanrio)",
		Link:  &feeds.Link{Href: s.baseURL + endpoint},
		Items: items,
	}

	return feed, nil
}
