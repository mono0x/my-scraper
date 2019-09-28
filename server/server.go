package server

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/feeds"
	"github.com/mono0x/my-scraper/scraper"
	"github.com/mono0x/my-scraper/scraper/source/fukokulifeevent"
	"github.com/mono0x/my-scraper/scraper/source/googlecalendar"
	"github.com/mono0x/my-scraper/scraper/source/harmonylandinfo"
	"github.com/mono0x/my-scraper/scraper/source/kittychaninfo"
	"github.com/mono0x/my-scraper/scraper/source/prtimes"
	"github.com/mono0x/my-scraper/scraper/source/purolandinfo"
	"github.com/mono0x/my-scraper/scraper/source/sanrionewsrelease"
	"github.com/mono0x/my-scraper/scraper/source/seibuenevent"
	"github.com/mono0x/my-scraper/scraper/source/yuyakekoyakenews"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

const cacheSeconds = 3600

func renderFeed(w http.ResponseWriter, feed *feeds.Feed) {
	w.Header().Set("Content-Type", "application/atom+xml")
	w.Header().Set("Cache-Control", "public, max-age="+string(cacheSeconds))
	if err := feed.WriteAtom(w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func sourceRenderer(source scraper.Source) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		feed, err := source.Scrape()
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
	}
}

func NewHandler() (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	entries := []struct {
		Path   string
		Source scraper.Source
	}{
		{"/fukoku-life", fukokulifeevent.NewSource(client)},
		{"/harmonyland-info", harmonylandinfo.NewSource(client)},
		{"/kittychan-info", kittychaninfo.NewSource(client)},
		{"/prtimes-sanrio", prtimes.NewSource(client)},
		{"/puroland-info", purolandinfo.NewSource(client)},
		{"/sanrio-news-release", sanrionewsrelease.NewSource(client)},
		{"/seibuen-event", seibuenevent.NewSource(client)},
		{"/yuyakekoyake-news", yuyakekoyakenews.NewSource(client)},
	}
	for _, entry := range entries {
		r.Get(entry.Path, sourceRenderer(entry.Source))
	}

	r.Get("/google-calendar", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := googlecalendar.NewSource(client, id)
		sourceRenderer(source)(w, r)
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
