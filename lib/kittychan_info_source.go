package scraper

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
	TitlePrefix      = `★`
)

var (
	titleDateRe = regexp.MustCompile(
		`\A` + TitlePrefix + `?(.+?)\s*(?:（(\d{4})年(\d{1,2})月(\d{1,2})日）)?\z`)
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

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	var items []*feeds.Item
	doc.Find("p").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		var (
			title, description, link string
		)
		s.Find("font").Each(func(_ int, s *goquery.Selection) {
			color, ok := s.Attr("color")
			if !ok {
				return
			}
			switch color {
			case "#0000ff":
				title = strings.TrimSpace(s.Text())
				break
			case "#000000":
				description, _ = s.Html()
				s.Find("a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
					if href, ok := s.Attr("href"); ok {
						link = href
						return false
					}
					return true
				})
				break
			}
		})

		matches := titleDateRe.FindStringSubmatch(title)
		if len(matches) < 2 || matches[1] == "" {
			return true
		}
		title = matches[1]

		var updated time.Time
		if len(matches) >= 5 && matches[2] != "" && matches[3] != "" && matches[4] != "" {
			year, err := strconv.Atoi(matches[2])
			if err != nil {
				return true
			}
			month, err := strconv.Atoi(matches[3])
			if err != nil {
				return true
			}
			day, err := strconv.Atoi(matches[4])
			if err != nil {
				return true
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
			if len(items) >= 100 {
				return false
			}

		}
		return true
	})
	feed.Items = items

	return feed, nil
}
