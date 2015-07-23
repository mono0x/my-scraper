package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"strings"
)

func GetCharacterShow() (*feeds.Feed, error) {
	doc, err := goquery.NewDocument("http://charactershow.jp/")
	if err != nil {
		return nil, err
	}
	return GetCharacterShowFromDocument(doc)
}

func GetCharacterShowFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title: "キャラクターショーファンサイト",
		Link:  &feeds.Link{Href: "http://charactershow.jp/"},
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
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(date+title+spot)))
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
