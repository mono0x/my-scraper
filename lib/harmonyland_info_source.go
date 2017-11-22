package scraper

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
)

const (
	harmonylandInfoURL = "http://www.harmonyland.jp/welcome.html"
)

type HarmonylandInfoSource struct {
}

var (
	harmonylandInfoTitleReplacer = strings.NewReplacer("\n", " ")
)

func NewHarmonylandInfoSource() *HarmonylandInfoSource {
	return &HarmonylandInfoSource{}
}

func (s *HarmonylandInfoSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get(harmonylandInfoURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	return s.ScrapeFromReader(res.Body)
}

func (s *HarmonylandInfoSource) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	decodedReader := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.ScrapeFromDocument(doc)
}

func (s *HarmonylandInfoSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	baseURL, _ := url.Parse(harmonylandInfoURL)

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

			resolvedHref := baseURL.ResolveReference(hrefURL).String()

			title := harmonylandInfoTitleReplacer.Replace(strings.TrimSpace(s.Text()))

			sha := sha256.New()
			fmt.Fprint(sha, title, resolvedHref)

			items = append(items, &feeds.Item{
				Title: title,
				Link:  &feeds.Link{Href: resolvedHref},
				Id:    fmt.Sprintf("%x", sha.Sum(nil)),
			})
		})
	})

	feed := &feeds.Feed{
		Title: "ハーモニーランド",
		Link:  &feeds.Link{Href: harmonylandInfoURL},
		Items: items,
	}
	return feed, nil
}
