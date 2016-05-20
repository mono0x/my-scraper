package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	MemoirsOfShibasakiSakiUrl = "http://shibasakisaki.web.fc2.com/"
)

var (
	TitleRe = regexp.MustCompile(`^\d+/\d+`)
)

type MemoirsOfShibasakiSakiSource struct {
}

func NewMemoirsOfShibasakiSakiSource() *MemoirsOfShibasakiSakiSource {
	return &MemoirsOfShibasakiSakiSource{}
}

func (s *MemoirsOfShibasakiSakiSource) Scrape() (*feeds.Feed, error) {
	res, err := http.Get(MemoirsOfShibasakiSakiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return s.ScrapeFromReader(res.Body)
}

func (s *MemoirsOfShibasakiSakiSource) ScrapeFromReader(reader io.Reader) (*feeds.Feed, error) {
	decodedReader := transform.NewReader(reader, japanese.ShiftJIS.NewDecoder())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, err
	}
	return s.ScrapeFromDocument(doc)
}

func (s *MemoirsOfShibasakiSakiSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	var items []*feeds.Item
	doc.Find(`td[bgcolor="#330066"] font[size="+1"] > *`).Each(func(_ int, s *goquery.Selection) {
		if s.Is("br") {
			return
		}

		var href string
		if s.Is("a") {
			href = s.AttrOr("href", "")
		} else {
			href = MemoirsOfShibasakiSakiUrl
		}

		text := strings.TrimSpace(s.Text())
		if text == "" {
			return
		}
		if !TitleRe.MatchString(text) {
			return
		}

		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(text+href)))

		items = append(items, &feeds.Item{
			Title: text,
			Link:  &feeds.Link{Href: href},
			Id:    hash,
		})
	})

	feed := &feeds.Feed{
		Title: "柴崎さきの見聞録",
		Link:  &feeds.Link{Href: MemoirsOfShibasakiSakiUrl},
		Items: items,
	}
	return feed, nil
}
