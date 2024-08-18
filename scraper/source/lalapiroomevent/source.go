package lalapiroomevent

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/mono0x/my-scraper/scraper"
)

const (
	baseURL  = "http://www.lifecorp.jp"
	endpoint = "/topics/lalapiroom/"
)

var dateRe = regexp.MustCompile(`(\d{4})年(\d{1,2})月(\d{1,2})日`)

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

func (*source) Name() string {
	return "lalapi-room-event"
}

func (s *source) Scrape(ctx context.Context, query url.Values) (*feeds.Feed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	res, err := s.httpClient.Do(req)
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

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return s.ScrapeFromDocument(doc)
}

func (s *source) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	href, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var items []*feeds.Item
	doc.Find("._schedule-box dl").Each(func(_ int, s *goquery.Selection) {
		dateStr := s.Find("dt").Text()
		location := s.Find("dd").Text()

		m := dateRe.FindStringSubmatch(dateStr)
		if len(m) != 4 {
			return
		}

		d := dateRe.FindString(dateStr)
		if d == "" {
			return
		}

		t, err := time.ParseInLocation("2006年1月2日", d, loc)
		if err != nil {
			return
		}

		var suffix string
		if i := strings.IndexByte(location, '\n'); i >= 0 {
			suffix = location[:i]
		} else {
			suffix = location
		}

		h := *href
		// Append a text fragment
		h.Fragment = fmt.Sprintf(":~:text=%s,-%s", dateStr, suffix)

		items = append(items, &feeds.Item{
			Title:   fmt.Sprintf("%s %s", dateStr, location),
			Created: t,
			Link:    &feeds.Link{Href: h.String()},
		})
	})

	feed := &feeds.Feed{
		Title: "ララピーのお部屋 イベントスケジュール",
		Link:  &feeds.Link{Href: href.String()},
		Items: items,
	}

	return feed, nil
}
