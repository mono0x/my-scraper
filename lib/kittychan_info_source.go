package scraper

import (
	"io"
	"net/http"
	"regexp"
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
	titleRe = regexp.MustCompile(
		`\A(?:` + regexp.QuoteMeta(TitlePrefix) + `)?(.+?)\s*(?:（(\d{4}年\d{1,2}月\d{1,2}日.*）))?\z`)
	dateRe = regexp.MustCompile(`\d{4}年\d{1,2}月\d{1,2}日`)
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

	skippedHrCount := 0
	doc.Find("hr, p").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if s.Is("hr") {
			skippedHrCount += 1
			return true
		}
		if skippedHrCount < 2 {
			return true
		}

		var (
			title, description, extraInfo, link string
		)
		s.ChildrenFiltered("font").Each(func(_ int, s *goquery.Selection) {
			color, ok := s.Attr("color")
			if !ok {
				return
			}
			switch color {
			case "#0000ff":
				matches := titleRe.FindStringSubmatch(strings.TrimSpace(s.Text()))
				if len(matches) == 3 {
					title = matches[1]
					extraInfo = matches[2]
				} else if len(matches) == 2 {
					title = matches[1]
					extraInfo = ""
				} else {
					title = ""
					extraInfo = ""
				}
				break
			case "#000000":
				{
					extraElement := s.Find(`b font[color="#ff0000"]`)
					if extraElement.Length() > 0 {
						extraInfo = strings.TrimSpace(extraElement.Text())
						break
					}
				}
				{
					description, _ = s.Html()
					s.Find("a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
						if href, ok := s.Attr("href"); ok {
							link = href
							return false
						}
						return true
					})
				}
				break
			}
		})

		var created time.Time
		if extraInfo != "" {
			matches := dateRe.FindStringSubmatch(extraInfo)
			if len(matches) > 0 {
				c, err := time.ParseInLocation("2006年1月2日", matches[0], loc)
				if err != nil {
					return true
				}
				created = c
			}
		}

		if title != "" && description != "" && link != "" {
			items = append(items, &feeds.Item{
				Title:       title,
				Created:     created,
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
