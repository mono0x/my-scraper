package scraper

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

const (
	HarmonylandInfoUrl = "http://www.harmonyland.jp/welcome.html"
)

type HarmonylandInfoSource struct {
}

func NewHarmonylandInfoSource() *HarmonylandInfoSource {
	return &HarmonylandInfoSource{}
}

func (s *HarmonylandInfoSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get(HarmonylandInfoUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return s.ScrapeFromReader(res.Body)
}

func (s *HarmonylandInfoSource) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	decodedReader := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, err
	}
	return s.ScrapeFromDocument(doc)
}

func (s *HarmonylandInfoSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	baseUrl, _ := url.Parse(HarmonylandInfoUrl)

	titleReplacer := strings.NewReplacer("\n", " ")

	var items []*feeds.Item
	doc.Find("#pickup, #cp").Each(func(_ int, s *goquery.Selection) {
		s.Find("dd .pickup_date").Each(func(_ int, s *goquery.Selection) {
			date, err := time.ParseInLocation("2006.1/2", s.Text(), loc)
			if err != nil {
				return
			}

			link := s.Next().Find("a")

			href, ok := link.Attr("href")
			if !ok {
				return
			}
			hrefUrl, err := url.Parse(href)
			if err != nil {
				return
			}

			title := titleReplacer.Replace(strings.TrimSpace(link.Text()))

			items = append(items, &feeds.Item{
				Title:   title,
				Link:    &feeds.Link{Href: baseUrl.ResolveReference(hrefUrl).String()},
				Created: date,
			})
		})
	})

	feed := &feeds.Feed{
		Title: "ハーモニーランド",
		Link:  &feeds.Link{Href: HarmonylandInfoUrl},
		Items: items,
	}
	return feed, nil
}
