package scraper

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

const (
	CharacterShowUrl = "http://charactershow.jp/"
)

type CharacterShowSource struct {
}

func NewCharacterShowSource() *CharacterShowSource {
	return &CharacterShowSource{}
}

func (s *CharacterShowSource) Scrape() (*feeds.Feed, error) {
	doc, err := goquery.NewDocument(CharacterShowUrl)
	if err != nil {
		return nil, err
	}
	return s.ScrapeFromDocument(doc)
}

func (s *CharacterShowSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title: "キャラクターショーファンサイト",
		Link:  &feeds.Link{Href: CharacterShowUrl},
	}

	dateReplacer := strings.NewReplacer("開催予定日", "", " ", "")
	titleReplacer := strings.NewReplacer("イベント内容", "")
	spotReplacer := strings.NewReplacer("開催場所(出典)", "")

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
		date := dateReplacer.Replace(s.Find(".event-date").Text())
		title := titleReplacer.Replace(s.Find(".event-name").Text())
		spot := spotReplacer.Replace(s.Find(".event-spotname").Text())
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
