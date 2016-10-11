package scraper

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

const (
	FukokuLifeEventUrl = "https://act.fukoku-life.co.jp/event/index.php"
)

type FukokuLifeEventSource struct {
}

func NewFukokuLifeEventSource() *FukokuLifeEventSource {
	return &FukokuLifeEventSource{}
}

func (s *FukokuLifeEventSource) Scrape() (*feeds.Feed, error) {
	doc, err := goquery.NewDocument(FukokuLifeEventUrl)
	if err != nil {
		return nil, err
	}
	return s.ScrapeFromDocument(doc)
}

func (s *FukokuLifeEventSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	baseUrl, _ := url.Parse(FukokuLifeEventUrl)
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
		refUrl, err := url.Parse(linkText)
		if err != nil {
			return
		}
		absUrl := baseUrl.ResolveReference(refUrl)
		link := absUrl.String()

		description := dateCell.Text() + "\n" + locationCell.Text()

		items = append(items, &feeds.Item{
			Title:       title,
			Description: description,
			Link:        &feeds.Link{Href: link},
			Id:          link,
		})
	})

	feed := &feeds.Feed{
		Title: "フコク赤ちゃんクラブ",
		Link:  &feeds.Link{Href: FukokuLifeEventUrl},
		Items: items,
	}

	return feed, nil
}
