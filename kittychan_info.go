package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

const (
	Url = "http://www.kittychan.info/information.html"
)

func GetKittychanInfo() (*feeds.Feed, error) {
	res, err := http.Get(Url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return GetKittychanInfoFromReader(res.Body)
}

func GetKittychanInfoFromReader(reader io.Reader) (*feeds.Feed, error) {
	cmd := exec.Command("iconv", "-c", "-f", "cp932", "-t", "utf-8")
	cmd.Stdin = reader
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	decodedReader := bytes.NewReader(out.Bytes())
	doc, err := goquery.NewDocumentFromReader(decodedReader)
	if err != nil {
		return nil, err
	}
	return GetKittychanInfoFromDocument(doc)
}

func GetKittychanInfoFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	feed := &feeds.Feed{
		Title: "♪キティちゃん情報",
		Link:  &feeds.Link{Href: Url},
	}

	var items []*feeds.Item
	doc.Find("p").Each(func(_ int, s *goquery.Selection) {
		if len(items) >= 100 {
			return
		}

		var title string
		var description string
		var link string
		s.Find("font").Each(func(_ int, s *goquery.Selection) {
			color, ok := s.Attr("color")
			if !ok {
				return
			}
			if color == "#0000ff" {
				title = strings.TrimSpace(s.Text())
			} else if color == "#000000" {
				description, _ = s.Html()
				s.Find("a").Each(func(_ int, s *goquery.Selection) {
					if link != "" {
						return
					}
					if href, ok := s.Attr("href"); ok {
						link = href
					}
				})
			}
		})

		if title != "" && description != "" && link != "" {
			items = append(items, &feeds.Item{
				Title:       title,
				Description: description,
				Link:        &feeds.Link{Href: link},
				Id:          link,
			})
		}
	})
	feed.Items = items

	return feed, nil
}
