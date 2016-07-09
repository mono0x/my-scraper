package scraper

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

const (
	SanrioNewsReleaseUrl = "http://www.sanrio.co.jp/corporate/release/"
)

type SanrioNewsReleaseSource struct {
}

func NewSanrioNewsReleaseSource() *SanrioNewsReleaseSource {
	return &SanrioNewsReleaseSource{}
}

func (s *SanrioNewsReleaseSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get(SanrioNewsReleaseUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return s.ScrapeFromDocument(doc)
}

func (s *SanrioNewsReleaseSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	var items []*feeds.Item

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
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
			Link:    &feeds.Link{Href: "http://www.sanrio.co.jp" + href},
			Id:      href,
			Created: date,
		})
	})

	feed := &feeds.Feed{
		Title: "ニュースリリース | 会社情報 | サンリオ",
		Link:  &feeds.Link{Href: SanrioNewsReleaseUrl},
		Items: items,
	}
	return feed, nil
}
