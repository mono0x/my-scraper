package yuyakekoyakenews

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/mono0x/my-scraper/scraper"
)

const (
	baseURL  = "http://yuyakekoyake.jp"
	endpoint = "/news/index.php"
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

func (s *source) Scrape(url.Values) (*feeds.Feed, error) {
	res, err := s.httpClient.Get(s.baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return s.ScrapeFromDocument(doc)
}

var yuyakekoyakeNewsItemRe = regexp.MustCompile(`\A(\d+)年(\d+)月(\d+)日[\s　]+(.+)\z`)

func (s *source) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	absBaseURL, _ := url.Parse(baseURL + endpoint)

	var items []*feeds.Item
	doc.Find(".news_detail_index li").Each(func(_ int, s *goquery.Selection) {
		m := yuyakekoyakeNewsItemRe.FindStringSubmatch(s.Text())
		if len(m) != 5 {
			return
		}

		year, err := strconv.Atoi(m[1])
		if err != nil {
			return
		}

		month, err := strconv.Atoi(m[2])
		if err != nil {
			return
		}

		day, err := strconv.Atoi(m[3])
		if err != nil {
			return
		}

		title := m[4]

		created := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)

		href, ok := s.Find("a").First().Attr("href")
		if !ok {
			return
		}

		refURL, err := url.Parse(href)
		if err != nil {
			return
		}

		absURL := absBaseURL.ResolveReference(refURL)

		items = append(items, &feeds.Item{
			Title:   title,
			Link:    &feeds.Link{Href: absURL.String()},
			Created: created,
		})
	})

	return &feeds.Feed{
		Title: "お知らせ一覧 | 夕やけ小やけふれあいの里",
		Link:  &feeds.Link{Href: baseURL + endpoint},
		Items: items,
	}, nil
}
