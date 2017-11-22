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
	seibuenEventURL = "http://www.seibuen-yuuenchi.jp/event/index.html?category=e1"
)

var (
	seibuenEventTitleReplacer = strings.NewReplacer("『", "", "』", "")
	seibuenEventTextReplacer  = strings.NewReplacer("\n", "", "\t", "")
)

type SeibuenEventSource struct {
}

func NewSeibuenEventSource() *SeibuenEventSource {
	return &SeibuenEventSource{}
}

func (s *SeibuenEventSource) Scrape() (*feeds.Feed, error) {
	doc, err := goquery.NewDocument(seibuenEventURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.ScrapeFromDocument(doc)
}

func (s *SeibuenEventSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	var items []*feeds.Item
	var (
		title string
	)
	doc.Find(".elem-section > div > div > div > div > div").Each(func(_ int, s *goquery.Selection) {
		switch {
		case s.HasClass("elem-heading-lv3"):
			title = seibuenEventTitleReplacer.Replace(s.Find("h3").Text())
		case s.HasClass("elem-pic-block"):
			properties := map[string]string{}
			s.Find("table tr").Each(func(_ int, s *goquery.Selection) {
				key := seibuenEventTextReplacer.Replace(strings.TrimSpace(s.Find("th").Text()))
				value := seibuenEventTextReplacer.Replace((strings.TrimSpace(s.Find("td").Text())))
				properties[key] = value
			})
			if len(properties) == 0 {
				return
			}

			summary := seibuenEventTextReplacer.Replace(s.Find(".txt-box > div > .txt-body > div > .elem-paragraph").Text())

			description := fmt.Sprintf("%s<br /><br />日程: %s<br />時間: %s<br />場所: %s<br />その他: %s", summary, properties["日程"], properties["時間"], properties["場所"], properties["その他"])

			sha := sha256.New()
			fmt.Fprint(sha, title, properties["日程"])

			items = append(items, &feeds.Item{
				Title:       title,
				Description: description,
				Link:        &feeds.Link{Href: seibuenEventURL},
				Id:          fmt.Sprintf("%x", sha.Sum(nil)),
			})
		}
	})

	feed := &feeds.Feed{
		Title: "西武園ゆうえんち メルヘンタウン",
		Link:  &feeds.Link{Href: seibuenEventURL},
		Items: items,
	}

	return feed, nil
}
