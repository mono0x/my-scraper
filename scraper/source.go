package scraper

import (
	"context"
	"net/url"

	"github.com/gorilla/feeds"
)

type Source interface {
	Name() string
	Scrape(ctx context.Context, query url.Values) (*feeds.Feed, error)
}
