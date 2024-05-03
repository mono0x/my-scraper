package kittychaninfo

import (
	"fmt"
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

func (s *source) Scrape(url.Values) (*feeds.Feed, error) {
	res, err := s.httpClient.Get(s.baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return &feeds.Feed{}, nil
		}
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return s.ScrapeFromReader(res.Body)
}

func (s *source) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	decodedReader := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return s.ScrapeFromDocument(doc)
}

func (s *source) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title: "♪キティちゃん情報",
		Link:  &feeds.Link{Href: s.baseURL + endpoint},
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
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
			var matches []string
			if createdStr != "" {
				matches = dateRe.FindStringSubmatch(createdStr)
			} else {
				// fallback to title or description
				matches = dateRe.FindStringSubmatch(title)
				if len(matches) == 0 {
					matches = dateRe.FindStringSubmatch(descriptionStr)
				}
			}
			if len(matches) > 0 {
				c, err := time.ParseInLocation("2006年1月2日", matches[0], loc)
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
