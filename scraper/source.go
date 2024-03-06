package scraper

import (
	"net/url"

	"github.com/gorilla/feeds"
)

type Source interface {
	Scrape(query url.Values) (*feeds.Feed, error)
}
