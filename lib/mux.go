package scraper

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/feeds"
	"github.com/mono0x/my-scraper/lib/facebook"
	"github.com/mono0x/my-scraper/lib/fukokulifeevent"
	"github.com/mono0x/my-scraper/lib/googlecalendar"
	"github.com/mono0x/my-scraper/lib/harmonylandinfo"
	"github.com/mono0x/my-scraper/lib/instagram"
	"github.com/mono0x/my-scraper/lib/kittychaninfo"
	"github.com/mono0x/my-scraper/lib/prtimes"
	"github.com/mono0x/my-scraper/lib/purolandinfo"
	"github.com/mono0x/my-scraper/lib/sanrionewsrelease"
	"github.com/mono0x/my-scraper/lib/seibuenevent"
	"github.com/mono0x/my-scraper/lib/twitter"
	"github.com/mono0x/my-scraper/lib/valuepress"
	"github.com/mono0x/my-scraper/lib/yuyakekoyakenews"
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

func sourceRenderer(source Source) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		feed, err := source.Scrape()
		if err != nil {
			log.Printf("%+v\n", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		renderFeed(w, feed)
	}
}

func NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	entries := []struct {
		Path   string
		Source Source
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
