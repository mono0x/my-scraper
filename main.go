package main

import (
  "fmt"
  "github.com/gorilla/feeds"
  "net/http"
)

func RenderFeed(w http.ResponseWriter, feed *feeds.Feed) {
  atom, err := feed.ToAtom()
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/atom+xml")
  fmt.Fprintln(w, atom)
}

func HandlePurolandNews(w http.ResponseWriter, r *http.Request) {
  feed, err := GetPurolandNews()
  if err != nil {
    w.WriteHeader(http.StatusServiceUnavailable)
    return
  }
  RenderFeed(w, feed)
}

func HandlePurolandInfo(w http.ResponseWriter, r *http.Request) {
  feed, err := GetPurolandInfo()
  if err != nil {
    w.WriteHeader(http.StatusServiceUnavailable)
    return
  }
  RenderFeed(w, feed)
}

func main() {
  http.HandleFunc("/puroland-info", HandlePurolandInfo)
  http.HandleFunc("/puroland-news", HandlePurolandNews)
  http.ListenAndServe(":13000", nil)
}
