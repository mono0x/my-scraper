package scraper

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
)

const (
	yuyakekoyakeNewsURL = "http://yuyakekoyake.jp/news/index.php"
)

type YuyakekoyakeNewsSource struct {
}

func NewYuyakekoyakeNewsSource() *YuyakekoyakeNewsSource {
	return &YuyakekoyakeNewsSource{}
}

func (s *YuyakekoyakeNewsSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get(yuyakekoyakeNewsURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.ScrapeFromDocument(doc)
}

var yuyakekoyakeNewsItemRe = regexp.MustCompile(`\A(\d+)年(\d+)月(\d+)日[\s　]+(.+)\z`)

func (s *YuyakekoyakeNewsSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	baseURL, _ := url.Parse(yuyakekoyakeNewsURL)

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

		absURL := baseURL.ResolveReference(refURL)

		items = append(items, &feeds.Item{
			Title:   title,
			Link:    &feeds.Link{Href: absURL.String()},
			Created: created,
		})
	})

	return &feeds.Feed{
		Title: "お知らせ一覧 | 夕やけ小やけふれあいの里",
		Link:  &feeds.Link{Href: yuyakekoyakeNewsURL},
		Items: items,
	}, nil
}
