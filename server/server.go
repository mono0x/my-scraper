package server

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/feeds"
	"github.com/mono0x/my-scraper/scraper"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

const cacheSeconds = 3600

var cacheControl = fmt.Sprintf("public, max-age=%d", cacheSeconds)

func renderFeed(w http.ResponseWriter, feed *feeds.Feed) {
	w.Header().Set("Content-Type", "application/atom+xml")
	w.Header().Set("Cache-Control", cacheControl)
	if err := feed.WriteAtom(w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func NewHandler(sources []scraper.Source) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	sourcesByName := make(map[string]scraper.Source)
	for _, source := range sources {
		sourcesByName[source.Name()] = source
	}

	r.Get("/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		source, ok := sourcesByName[name]
		if !ok {
			log.Printf("%v: not found\n", name)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		feed, err := source.Scrape(r.URL.Query())
		if err != nil {
			log.Printf("%v: %+v\n", reflect.TypeOf(source), err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		if len(feed.Items) == 0 {
			log.Printf("%v: not found\n", reflect.TypeOf(source))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		renderFeed(w, feed)
	})

	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(1024),
	)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(cacheSeconds*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return cacheClient.Middleware(r), nil
}
