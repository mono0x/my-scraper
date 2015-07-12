package main

import (
  "fmt"
  "github.com/gorilla/feeds"
  "github.com/PuerkitoBio/goquery"
  "net/http"
  "strings"
)

func handleFeed(w http.ResponseWriter, feed *feeds.Feed) {
  atom, err := feed.ToAtom()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/atom+xml")
  fmt.Fprintln(w, atom)
}

func getPurolandNews() (*feeds.Feed, error) {
  doc, err := goquery.NewDocument("http://www.puroland.jp/")
  if err != nil {
    return nil, err
  }

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

func handlePurolandNews(w http.ResponseWriter, r *http.Request) {
  feed, err := getPurolandNews()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  handleFeed(w, feed)
}

func getPurolandInfo() (*feeds.Feed, error) {
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

func handlePurolandInfo(w http.ResponseWriter, r *http.Request) {
  feed, err := getPurolandInfo()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  handleFeed(w, feed)
}

func main() {
  http.HandleFunc("/puroland-info", handlePurolandInfo)
  http.HandleFunc("/puroland-news", handlePurolandNews)
  http.ListenAndServe(":13000", nil)
}
