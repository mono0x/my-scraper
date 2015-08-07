package main

import (
	//"crypto/sha256"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"strings"
)

func GetSanrioEvent() (*feeds.Feed, error) {
	doc, err := goquery.NewDocument("http://www.sanrio.co.jp/event/search/")
	if err != nil {
		return nil, err
	}
	return GetSanrioEventFromDocument(doc)
}

func GetSanrioEventFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title: "イベント情報 | サンリオ",
		Link:  &feeds.Link{Href: "http://www.sanrio.co.jp/"},
	}

	spotReplacer := strings.NewReplacer("会場：", "")
	dateReplacer := strings.NewReplacer("開催日：", "")

	var items []*feeds.Item
	doc.Find(".ev_box").Each(func(_ int, s *goquery.Selection) {
		title := s.Find(".region_name").Text()
		spot := spotReplacer.Replace(s.Find(".ev_place").Text())
		date := dateReplacer.Replace(s.Find(".ev_schedule *:first-child").Text())
		link, ok := s.Find(".region_name a").Attr("href")
		if !ok {
			return
		}
		items = append(items, &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: link},
			Id:          link,
			Description: fmt.Sprintf("%s: %s", date, spot),
		})
	})
	feed.Items = items

	return feed, nil
}
