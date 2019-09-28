package fukokulifeevent

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/mono0x/my-scraper/scraper"
)

const (
	baseURL  = "https://act.fukoku-life.co.jp"
	endpoint = "/event/index.php"
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
	absBaseURL, _ := url.Parse(baseURL + endpoint)
	var items []*feeds.Item
	doc.Find("div#result > table > tbody > tr").Each(func(_ int, s *goquery.Selection) {
		titleCell := s.Children().First().Next()
		dateCell := titleCell.Next()
		locationCell := dateCell.Next()

		title := titleCell.Text()
		linkText, ok := titleCell.Children().First().Attr("href")
		if !ok {
			return
		}
		refURL, err := url.Parse(linkText)
		if err != nil {
			return
		}
		absURL := absBaseURL.ResolveReference(refURL)
		link := absURL.String()

		description := dateCell.Text() + "\n" + locationCell.Text()

		items = append(items, &feeds.Item{
			Title:       title,
			Description: description,
			Link:        &feeds.Link{Href: link},
		})
	})

	feed := &feeds.Feed{
		Title: "フコク赤ちゃんクラブ",
		Link:  &feeds.Link{Href: baseURL + endpoint},
		Items: items,
	}

	return feed, nil
}
