package scraper

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

const (
	PRTimesUrl = `http://prtimes.jp/main/action.php?run=html&page=searchkey&search_word=%E3%82%B5%E3%83%B3%E3%83%AA%E3%82%AA&search_pattern=1`
)

type PRTimesSource struct {
}

func NewPRTimesSource() *PRTimesSource {
	return &PRTimesSource{}
}

func (s *PRTimesSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get(PRTimesUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return s.ScrapeFromReader(res.Body)
}

func (s *PRTimesSource) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	return s.ScrapeFromDocument(doc)
}

func (s *PRTimesSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title: "PR TIMES (Sanrio)",
		Link:  &feeds.Link{Href: PRTimesUrl},
	}

	baseUrl, _ := url.Parse(PRTimesUrl)
	var items []*feeds.Item
	doc.Find("a.link-title-item-ordinary").Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		link, ok := s.Attr("href")
		if !ok {
			return
		}
		refUrl, err := url.Parse(link)
		if err != nil {
			return
		}
		absUrl := baseUrl.ResolveReference(refUrl)
		link = absUrl.String()
		dt, ok := s.Parent().Next().Attr("datetime")
		if !ok {
			return
		}
		t, err := time.Parse("2006-01-02T15:04:05-0700", dt)
		if err != nil {
			return
		}
		items = append(items, &feeds.Item{
			Title:   title,
			Link:    &feeds.Link{Href: link},
			Id:      link,
			Created: t,
		})
	})
	feed.Items = items

	return feed, nil
}
