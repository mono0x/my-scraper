package scraper

import "github.com/gorilla/feeds"

type Source interface {
	Scrape() (*feeds.Feed, error)
}
