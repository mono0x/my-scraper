package lalapiroomevent

import (
	"context"
	"errors"
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

var (
	dateRe = regexp.MustCompile(`\d{4}年\d{1,2}月\d{1,2}日`)
	jst    = time.FixedZone("JST", 9*60*60)
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

func (*source) Name() string {
	return "lalapi-room-event"
}

func (s *source) Scrape(ctx context.Context, query url.Values) (*feeds.Feed, error) {
	body, err := scraper.Fetch(ctx, s.httpClient, s.baseURL+endpoint)
	if err != nil {
		if errors.Is(err, scraper.ErrNotFound) {
			return &feeds.Feed{}, nil
		}
		return nil, err
	}
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	return s.ScrapeFromDocument(doc)
}

func (s *source) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	href, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	var items []*feeds.Item
	doc.Find("._schedule-box dl").Each(func(_ int, s *goquery.Selection) {
		dateStr := s.Find("dt").Text()
		location := s.Find("dd").Text()

		d := dateRe.FindString(dateStr)
		if d == "" {
			return
		}

		t, err := time.ParseInLocation("2006年1月2日", d, jst)
		if err != nil {
			return
		}

		var suffix string
		if before, _, ok := strings.Cut(location, "\n"); ok {
			suffix = before
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
