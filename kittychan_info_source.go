package main

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

const (
	KittychanInfoUrl = "http://www.kittychan.info/information.html"
)

type KittychanInfoSource struct {
}

func NewKittychanInfoSource() *KittychanInfoSource {
	return &KittychanInfoSource{}
}

func (s *KittychanInfoSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get(KittychanInfoUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return s.ScrapeFromReader(res.Body)
}

func (s *KittychanInfoSource) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	decodedReader := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, err
	}
	return s.ScrapeFromDocument(doc)
}

func (s *KittychanInfoSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title: "♪キティちゃん情報",
		Link:  &feeds.Link{Href: KittychanInfoUrl},
	}

	titleDateRe := regexp.MustCompile(`\A★?(.+?)\s*(?:（(\d{4})年(\d{1,2})月(\d{1,2})日）)?\z`)

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	var items []*feeds.Item
	doc.Find("p").Each(func(_ int, s *goquery.Selection) {
		if len(items) >= 100 {
			return
		}

		var title string
		var description string
		var link string
		s.Find("font").Each(func(_ int, s *goquery.Selection) {
			color, ok := s.Attr("color")
			if !ok {
				return
			}
			if color == "#0000ff" {
				title = strings.TrimSpace(s.Text())
			} else if color == "#000000" {
				description, _ = s.Html()
				s.Find("a").Each(func(_ int, s *goquery.Selection) {
					if link != "" {
						return
					}
					if href, ok := s.Attr("href"); ok {
						link = href
					}
				})
			}
		})

		matches := titleDateRe.FindStringSubmatch(title)
		if len(matches) < 2 || matches[1] == "" {
			return
		}
		title = matches[1]

		var updated time.Time
		if len(matches) >= 5 && matches[2] != "" && matches[3] != "" && matches[4] != "" {
			year, err := strconv.Atoi(matches[2])
			if err != nil {
				return
			}
			month, err := strconv.Atoi(matches[3])
			if err != nil {
				return
			}
			day, err := strconv.Atoi(matches[4])
			if err != nil {
				return
			}
			updated = time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
		}

		if title != "" && description != "" && link != "" {
			items = append(items, &feeds.Item{
				Title:       title,
				Updated:     updated,
				Description: description,
				Link:        &feeds.Link{Href: link},
				Id:          link,
			})
		}
	})
	feed.Items = items

	return feed, nil
}
