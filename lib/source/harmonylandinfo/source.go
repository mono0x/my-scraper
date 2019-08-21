package harmonylandinfo

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/pkg/errors"
)

const (
	baseURL  = "http://www.harmonyland.jp"
	endpoint = "/welcome.html"
)

var (
	titleReplacer = strings.NewReplacer("\n", " ")
)

type source struct {
	httpClient *http.Client
	baseURL    string
}

var _ scraper.Source = (*source)(nil)

func NewSource(c *http.Client) *source {
	return &source{
		httpClient: c,
		baseURL:    baseURL,
	}
}

func (s *source) Scrape() (*feeds.Feed, error) {
	res, err := s.httpClient.Get(s.baseURL + endpoint)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	return s.scrapeFromReader(res.Body)
}

func (s *source) scrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	decodedReader := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.scrapeFromDocument(doc)
}

func (s *source) scrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	absBaseURL, _ := url.Parse(baseURL + endpoint)

	var items []*feeds.Item
	doc.Find("#pickup, #cp").Each(func(_ int, s *goquery.Selection) {
		s.Find("dd .pick_up").Each(func(_ int, s *goquery.Selection) {
			link := s.Find("a")

			href, ok := link.Attr("href")
			if !ok {
				return
			}
			hrefURL, err := url.Parse(href)
			if err != nil {
				return
			}

			resolvedHref := absBaseURL.ResolveReference(hrefURL).String()

			title := titleReplacer.Replace(strings.TrimSpace(s.Text()))

			items = append(items, &feeds.Item{
				Title: title,
				Link:  &feeds.Link{Href: resolvedHref},
			})
		})
	})

	feed := &feeds.Feed{
		Title: "ハーモニーランド",
		Link:  &feeds.Link{Href: s.baseURL + endpoint},
		Items: items,
	}
	return feed, nil
}
