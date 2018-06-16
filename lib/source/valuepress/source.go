package valuepress

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
)

const (
	valuePressURL = `https://www.value-press.com/search?q=%E3%82%B5%E3%83%B3%E3%83%AA%E3%82%AA`
)

type ValuePressSource struct {
}

func NewSource() *ValuePressSource {
	return &ValuePressSource{}
}

func (s *ValuePressSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get(valuePressURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	return s.ScrapeFromReader(res.Body)
}

func (s *ValuePressSource) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.ScrapeFromDocument(doc)
}

func (s *ValuePressSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
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
		Link:  &feeds.Link{Href: valuePressURL},
		Items: items,
	}

	return feed, nil
}
