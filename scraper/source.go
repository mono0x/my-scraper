package scraper

import (
	"net/url"

	"github.com/gorilla/feeds"
)

type Source interface {
	Name() string
	Scrape(query url.Values) (*feeds.Feed, error)
}
