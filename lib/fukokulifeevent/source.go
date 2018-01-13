package fukokulifeevent

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
)

const (
	fukokuLifeEventURL = "https://act.fukoku-life.co.jp/event/index.php"
)

type FukokuLifeEventSource struct {
}

func NewSource() *FukokuLifeEventSource {
	return &FukokuLifeEventSource{}
}

func (s *FukokuLifeEventSource) Scrape() (*feeds.Feed, error) {
	doc, err := goquery.NewDocument(fukokuLifeEventURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.ScrapeFromDocument(doc)
}

func (s *FukokuLifeEventSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	baseURL, _ := url.Parse(fukokuLifeEventURL)
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
		absURL := baseURL.ResolveReference(refURL)
		link := absURL.String()

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
		Link:  &feeds.Link{Href: fukokuLifeEventURL},
		Items: items,
	}

	return feed, nil
}
