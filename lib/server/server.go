package server

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/feeds"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/mono0x/my-scraper/lib/source/facebook"
	"github.com/mono0x/my-scraper/lib/source/fukokulifeevent"
	"github.com/mono0x/my-scraper/lib/source/googlecalendar"
	"github.com/mono0x/my-scraper/lib/source/harmonylandinfo"
	"github.com/mono0x/my-scraper/lib/source/instagram"
	"github.com/mono0x/my-scraper/lib/source/kittychaninfo"
	"github.com/mono0x/my-scraper/lib/source/prtimes"
	"github.com/mono0x/my-scraper/lib/source/purolandinfo"
	"github.com/mono0x/my-scraper/lib/source/sanrionewsrelease"
	"github.com/mono0x/my-scraper/lib/source/seibuenevent"
	"github.com/mono0x/my-scraper/lib/source/twitter"
	"github.com/mono0x/my-scraper/lib/source/valuepress"
	"github.com/mono0x/my-scraper/lib/source/yuyakekoyakenews"
	"github.com/pkg/errors"
	"github.com/victorspringer/http-cache"
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

	entries := []struct {
		Path   string
		Source scraper.Source
	}{
		{"/fukoku-life", fukokulifeevent.NewSource(http.DefaultClient)},
		{"/harmonyland-info", harmonylandinfo.NewSource(http.DefaultClient)},
		{"/kittychan-info", kittychaninfo.NewSource(http.DefaultClient)},
		{"/prtimes-sanrio", prtimes.NewSource(http.DefaultClient)},
		{"/puroland-info", purolandinfo.NewSource(http.DefaultClient)},
		{"/sanrio-news-release", sanrionewsrelease.NewSource(http.DefaultClient)},
		{"/seibuen-event", seibuenevent.NewSource(http.DefaultClient)},
		{"/value-press-sanrio", valuepress.NewSource(http.DefaultClient)},
		{"/yuyakekoyake-news", yuyakekoyakenews.NewSource(http.DefaultClient)},
	}
	for _, entry := range entries {
		r.Get(entry.Path, sourceRenderer(entry.Source))
	}

	r.Get("/facebook", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := facebook.NewSource(http.DefaultClient, os.Getenv("FACEBOOK_ACCESS_TOKEN"), id)
		sourceRenderer(source)(w, r)
	})

	r.Get("/google-calendar", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := googlecalendar.NewSource(http.DefaultClient, id)
		sourceRenderer(source)(w, r)
	})

	r.Get("/instagram", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := instagram.NewSource(http.DefaultClient, id)
		sourceRenderer(source)(w, r)
	})

	r.Get("/twitter", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		idStr := query.Get("id")
		if idStr == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := twitter.NewSource(http.DefaultClient, id)
		sourceRenderer(source)(w, r)
	})

	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(1024),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(cacheSeconds*time.Second),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cacheClient.Middleware(r), nil
}
