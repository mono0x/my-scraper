package main

import (
  "github.com/gorilla/feeds"
  "github.com/PuerkitoBio/goquery"
  "strings"
)

func GetPurolandNews() (*feeds.Feed, error) {
  doc, err := goquery.NewDocument("http://www.puroland.jp/")
  if err != nil {
    return nil, err
  }
  return GetPurolandNewsFromDocument(doc)
}

func GetPurolandNewsFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
  feed := &feeds.Feed{
    Title: "最新情報 | サンリオピューロランド",
    Link: &feeds.Link{Href: "http://www.puroland.jp/"},
  }

  var items []*feeds.Item
  doc.Find("#newsArea ul li a").Each(func(_ int, s *goquery.Selection) {
    title := strings.TrimSpace(s.Text())
    link, ok := s.Attr("href")
    if ok {
      items = append(items, &feeds.Item{
        Title: title,
        Link: &feeds.Link{Href: link},
        Id: link,
      })
    }
  })
  feed.Items = items

  return feed, nil
}

func GetPurolandInfo() (*feeds.Feed, error) {
  doc, err := goquery.NewDocument("http://www.puroland.jp/")
  if err != nil {
    return nil, err
  }
  return GetPurolandInfoFromDocument(doc)
}

func GetPurolandInfoFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
  doc, err := goquery.NewDocument("http://www.puroland.jp/")
  if err != nil {
    return nil, err
  }

  feed := &feeds.Feed{
    Title: "お知らせ | サンリオピューロランド",
    Link: &feeds.Link{Href: "http://www.puroland.jp/"},
  }

  var items []*feeds.Item
  doc.Find("#infoSectionArea ul li a").Each(func(_ int, s *goquery.Selection) {
    title := strings.TrimSpace(s.Text())
    link, ok := s.Attr("href")
    if ok {
      items = append(items, &feeds.Item{
        Title: title,
        Link: &feeds.Link{Href: link},
        Id: link,
      })
    }
  })
  feed.Items = items

  return feed, nil
}
