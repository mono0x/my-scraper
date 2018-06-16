package server

import (
	"log"
	"net/http"
	"strconv"

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
)

func renderFeed(w http.ResponseWriter, feed *feeds.Feed) {
	w.Header().Set("Content-Type", "application/atom+xml")
	w.Header().Set("Cache-Control", "public, max-age=3600")
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
			log.Printf("%+v\n", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		if len(feed.Items) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		renderFeed(w, feed)
	}
}

func NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	entries := []struct {
		Path   string
		Source scraper.Source
	}{
		{"/fukoku-life", fukokulifeevent.NewSource()},
		{"/harmonyland-info", harmonylandinfo.NewSource()},
		{"/kittychan-info", kittychaninfo.NewSource()},
		{"/prtimes-sanrio", prtimes.NewSource()},
		{"/puroland-info", purolandinfo.NewSource()},
		{"/sanrio-news-release", sanrionewsrelease.NewSource()},
		{"/seibuen-event", seibuenevent.NewSource()},
		{"/value-press-sanrio", valuepress.NewSource()},
		{"/yuyakekoyake-news", yuyakekoyakenews.NewSource()},
	}
	for _, entry := range entries {
		mux.HandleFunc(entry.Path, sourceRenderer(entry.Source))
	}

	mux.HandleFunc("/facebook", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := facebook.NewSource(id)
		sourceRenderer(source)(w, r)
	})

	mux.HandleFunc("/google-calendar", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := googlecalendar.NewSource(id)
		sourceRenderer(source)(w, r)
	})

	mux.HandleFunc("/instagram", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := instagram.NewSource(id)
		sourceRenderer(source)(w, r)
	})

	mux.HandleFunc("/twitter", func(w http.ResponseWriter, r *http.Request) {
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
		source := twitter.NewSource(id)
		sourceRenderer(source)(w, r)
	})

	return mux
}
