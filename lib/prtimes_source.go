package scraper

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

const (
	PRTimesUrl = "http://prtimes.jp/topics/keywords/サンリオ"
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
		Title: "PRTimes",
		Link:  &feeds.Link{Href: PRTimesUrl},
	}

	var items []*feeds.Item
	doc.Find("a.link-title-item-ordinary").Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		link, ok := s.Attr("href")
		if !ok {
			return
		}
		dt, ok := s.Parent().Next().Next().Attr("datetime")
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
