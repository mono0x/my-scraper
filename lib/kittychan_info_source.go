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
	"github.com/pkg/errors"
)

const (
	kittychanInfoURL         = "http://www.kittychan.info/information.html"
	kittychanInfoTitlePrefix = `★`
)

var (
	headerRe = regexp.MustCompile(
		`\A(?:` + regexp.QuoteMeta(kittychanInfoTitlePrefix) + `)?(.+?)\s*(?:（(\d{4}年\d{1,2}月\d{1,2}日.*）))?\z`)
	dateRe = regexp.MustCompile(`\d{4}年\d{1,2}月\d{1,2}日`)

	descriptionReplacer = strings.NewReplacer("\n", "<br />")
)

type KittychanInfoSource struct {
}

func NewKittychanInfoSource() *KittychanInfoSource {
	return &KittychanInfoSource{}
}

func (s *KittychanInfoSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get(kittychanInfoURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	return s.ScrapeFromReader(res.Body)
}

func (s *KittychanInfoSource) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	decodedReader := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.ScrapeFromDocument(doc)
}

func (s *KittychanInfoSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	headerRe := headerRe.Copy()
	dateRe := dateRe.Copy()

	feed := &feeds.Feed{
		Title: "♪キティちゃん情報",
		Link:  &feeds.Link{Href: kittychanInfoURL},
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var items []*feeds.Item

	skippedHrCount := 0
	section := ""
	link := ""
	doc.Find("hr, p").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if s.Is("hr") {
			skippedHrCount += 1
			return true
		}
		if skippedHrCount < 2 {
			return true
		}

		p := strings.TrimSpace(s.Text())

		var extractedLink string
		s.Find("a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			if href, ok := s.Attr("href"); ok {
				extractedLink = href
				return false
			}
			return true
		})

		if !strings.HasPrefix(p, kittychanInfoTitlePrefix) {
			section += p
			if link == "" {
				link = extractedLink
			}
			return true
		}

		defer func() {
			section = p
			link = extractedLink
		}()

		parts := strings.SplitN(section, "\n", 2)
		if len(parts) < 2 {
			return true
		}

		header := parts[0]
		description := descriptionReplacer.Replace(parts[1])

		var title, extraInfo string
		{
			matches := headerRe.FindStringSubmatch(header)
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
		}

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
