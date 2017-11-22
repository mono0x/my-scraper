package scraper

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
)

const (
	characterShowURL = "http://charactershow.jp/"
)

var (
	characterShowDateReplacer  = strings.NewReplacer("開催予定日", "", " ", "")
	characterShowTitleReplacer = strings.NewReplacer("イベント内容", "")
	characterShowSpotReplacer  = strings.NewReplacer("開催場所(出典)", "")
)

type CharacterShowSource struct {
}

func NewCharacterShowSource() *CharacterShowSource {
	return &CharacterShowSource{}
}

func (s *CharacterShowSource) Scrape() (*feeds.Feed, error) {
	doc, err := goquery.NewDocument(characterShowURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.ScrapeFromDocument(doc)
}

func (s *CharacterShowSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title: "キャラクターショーファンサイト",
		Link:  &feeds.Link{Href: characterShowURL},
	}

	var items []*feeds.Item
	pref := ""
	prefTitle := ""
	doc.Find(".pref-title-jump, .pref-title, .event-row").Each(func(_ int, s *goquery.Selection) {
		if s.HasClass("pref-title-jump") {
			pref, _ = s.Attr("id")
			return
		}
		if s.HasClass("pref-title") {
			prefTitle = strings.TrimSpace(s.Find("div:first-child").Text())
			return
		}
		date := characterShowDateReplacer.Replace(s.Find(".event-date").Text())
		title := characterShowTitleReplacer.Replace(s.Find(".event-name").Text())
		spot := characterShowSpotReplacer.Replace(s.Find(".event-spotname").Text())
		link := "http://charactershow.jp/#" + pref
		sha := sha256.New()
		fmt.Fprint(sha, date, title, spot)
		hash := fmt.Sprintf("%x", sha.Sum(nil))
		items = append(items, &feeds.Item{
			Title:       fmt.Sprintf("%s: %s", prefTitle, title),
			Link:        &feeds.Link{Href: link},
			Id:          hash,
			Description: fmt.Sprintf("%s: %s", date, spot),
		})
	})
	feed.Items = items

	return feed, nil
}
