package impresswatchcolumn

import (
	"context"
	"errors"
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
	jst           = time.FixedZone("JST", 9*60*60)
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

	body, err := scraper.Fetch(ctx, s.httpClient, r.Replace(s.baseURL+endpoint))
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
	return s.ScrapeFromDocument(doc, r.Replace(baseURL+endpoint))
}

func (s *source) ScrapeFromDocument(doc *goquery.Document, siteURL string) (*feeds.Feed, error) {
	var items []*feeds.Item
	doc.Find("#main .list .list-02 .item").Each(func(_ int, s *goquery.Selection) {
		titleElement := s.Find(".title a")
		title := titleElement.Text()
		href, exists := titleElement.Attr("href")
		if !exists {
			return
		}

		dateStr := s.Find(".date").Text()
		t, err := time.ParseInLocation("(2006/1/2)", dateStr, jst)
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
