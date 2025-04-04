package impresswatchcolumn

import (
	"context"
	"fmt"
	"html"
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
	baseURL  = "https://{site}.watch.impress.co.jp"
	endpoint = "/docs/column/{column}/"
)

var (
	siteRe        = regexp.MustCompile(`\A[a-z0-9-]+\z`)
	titleSuffixRe = regexp.MustCompile(`\s+\d+年\s+記事一覧\z`)
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
	return "impress-watch-column"
}

func (s *source) Scrape(ctx context.Context, query url.Values) (*feeds.Feed, error) {
	site := query.Get("site")
	column := query.Get("column")
	if site == "" || column == "" {
		return &feeds.Feed{}, nil
	}
	if !siteRe.MatchString(site) {
		return nil, fmt.Errorf("invalid site: %s", site)
	}

	r := strings.NewReplacer("{site}", site, "{column}", column)

	req, err := http.NewRequestWithContext(ctx, "GET", r.Replace(s.baseURL+endpoint), nil)
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
	return s.ScrapeFromDocument(doc, r.Replace(baseURL+endpoint))
}

func (s *source) ScrapeFromDocument(doc *goquery.Document, siteURL string) (*feeds.Feed, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var items []*feeds.Item
	doc.Find("#main .list .list-02 .item").Each(func(_ int, s *goquery.Selection) {
		titleElement := s.Find(".title a")
		title := titleElement.Text()
		href, exists := titleElement.Attr("href")
		if !exists {
			return
		}

		dateStr := s.Find(".date").Text()
		t, err := time.ParseInLocation("(2006/1/2)", dateStr, loc)
		if err != nil {
			return
		}

		description := ""
		src, exists := s.Find(".image img").Attr("ajax") // src is not available in HTML due to lazy loading
		if exists {
			description += `<img src="` + html.EscapeString(src) + `" width="360" height="270" /><br />`
		}

		outline := s.Find(".outline").Text()
		if outline != "" {
			description += `<p>` + html.EscapeString(outline) + `</p>`
		}

		items = append(items, &feeds.Item{
			Title:       title,
			Description: description,
			Created:     t,
			Link:        &feeds.Link{Href: href},
		})
	})

	feed := &feeds.Feed{
		Title: titleSuffixRe.ReplaceAllString(doc.Find("#main article[role=main] > .title").Text(), ""),
		Link:  &feeds.Link{Href: siteURL},
		Items: items,
	}

	return feed, nil
}
