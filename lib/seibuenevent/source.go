package seibuenevent

import (
	"crypto/sha256"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
)

const (
	seibuenEventURL = "http://www.seibu-leisure.co.jp/event/index.html?category=e1"
)

type SeibuenEventSource struct {
}

func NewSource() *SeibuenEventSource {
	return &SeibuenEventSource{}
}

func (s *SeibuenEventSource) Scrape() (*feeds.Feed, error) {
	doc, err := goquery.NewDocument(seibuenEventURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.ScrapeFromDocument(doc)
}

func (s *SeibuenEventSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
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
				Link:        &feeds.Link{Href: seibuenEventURL},
				Id:          fmt.Sprintf("%x", sha.Sum(nil)),
			})
		}
	})

	feed := &feeds.Feed{
		Title: "西武園ゆうえんち メルヘンタウン",
		Link:  &feeds.Link{Href: seibuenEventURL},
		Items: items,
	}

	return feed, nil
}
