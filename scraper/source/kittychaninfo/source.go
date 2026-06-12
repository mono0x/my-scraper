package kittychaninfo

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/mono0x/my-scraper/scraper"
)

const (
	baseURL  = "http://www.kittychan.info"
	endpoint = "/information.html"

	titlePrefix = `★`
)

var (
	dateRe = regexp.MustCompile(`\d{4}年\d{1,2}月\d{1,2}日`)

	descriptionReplacer = strings.NewReplacer("\n", "<br />")

	jst = time.FixedZone("JST", 9*60*60)
)

type source struct {
	httpClient *http.Client
	baseURL    string // for testing
}

var _ scraper.Source = (*source)(nil)

func NewSource(c *http.Client) *source {
	return &source{
		httpClient: c,
		baseURL:    baseURL,
	}
}

func (s *source) Name() string {
	return "kittychan-info"
}

func (s *source) Scrape(ctx context.Context, _ url.Values) (*feeds.Feed, error) {
	body, err := scraper.Fetch(ctx, s.httpClient, s.baseURL+endpoint)
	if err != nil {
		if errors.Is(err, scraper.ErrNotFound) {
			return &feeds.Feed{}, nil
		}
		return nil, err
	}
	defer body.Close()

	return s.ScrapeFromReader(body)
}

func (s *source) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	decodedReader := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, err
	}
	return s.ScrapeFromDocument(doc)
}

func (s *source) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title: "♪キティちゃん情報",
		Link:  &feeds.Link{Href: s.baseURL + endpoint},
	}

	var items []*feeds.Item

	skippedHrCount := 0
	doc.Find("hr, p").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		// skip heading of the page
		if s.Is("hr") {
			skippedHrCount += 1
			return true
		}
		if skippedHrCount < 2 {
			return true
		}

		titleBlock := s.Find(`font[color="#0000ff"]`)
		title := strings.TrimPrefix(titleBlock.Text(), titlePrefix)
		createdStr := titleBlock.Find(`font[color="#ff0000"]`).Text()
		descriptionStr := s.Find(`font[color="#000000"]`).Text()

		var link string
		s.Find("a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			if href, ok := s.Attr("href"); ok {
				link = href
				return false
			}
			return true
		})

		var created time.Time
		{
			var d string
			if createdStr != "" {
				d = dateRe.FindString(createdStr)
			} else {
				// fallback to title or description
				d = dateRe.FindString(title)
				if d == "" {
					d = dateRe.FindString(descriptionStr)
				}
			}
			if d != "" {
				c, err := time.ParseInLocation("2006年1月2日", d, jst)
				if err != nil {
					return true
				}
				created = c
			}
		}

		if title != "" && descriptionStr != "" && link != "" {
			description := descriptionReplacer.Replace(descriptionStr)
			items = append(items, &feeds.Item{
				Title:       title,
				Created:     created,
				Description: description,
				Link:        &feeds.Link{Href: link},
			})
			if len(items) >= 100 {
				return false // break
			}

		}
		return true
	})
	feed.Items = items

	return feed, nil
}
