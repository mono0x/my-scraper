package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/feeds"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	PurolandInfoUrl    = "http://www.puroland.jp/"
	PurolandInfoApiUrl = "http://www.puroland.jp/api/live/get_information/?page=1&count=20"
)

type information struct {
	Status  string            `json:"status"`
	Count   int               `json:"count"`
	Total   int               `json:"total"`
	Page    int               `json:"page"`
	MaxPage int               `json:"maxpage"`
	Data    []informationItem `json:"data"`
}

type informationItem struct {
	Url             string `json:"url"`
	Title           string `json:"title"`
	PublicDate      string `json:"public_date"`
	ThumbnailMiddle string `json:"thumbnail_m"`
}

func GetPurolandInfo() (*feeds.Feed, error) {
	res, err := http.Get(PurolandInfoApiUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	feed, err := GetPurolandInfoFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func GetPurolandInfoFromReader(reader io.Reader) (*feeds.Feed, error) {
	jsonData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var info information
	err = json.Unmarshal(jsonData, &info)
	if err != nil {
		return nil, err
	}

	items := make([]*feeds.Item, info.Count)
	for i, infoItem := range info.Data {
		created, err := time.Parse("2006/01/02", infoItem.PublicDate)
		if err != nil {
			return nil, err
		}

		description := fmt.Sprintf("<img src=\"%s\" />", infoItem.ThumbnailMiddle)

		items[i] = &feeds.Item{
			Title:       html.UnescapeString(infoItem.Title),
			Link:        &feeds.Link{Href: infoItem.Url},
			Id:          infoItem.Url,
			Created:     created,
			Description: description,
		}
	}

	feed := &feeds.Feed{
		Title: "お知らせ | サンリオピューロランド",
		Link:  &feeds.Link{Href: PurolandInfoUrl},
		Items: items,
	}

	return feed, nil
}
