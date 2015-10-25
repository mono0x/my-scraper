package main

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func renderFeed(w http.ResponseWriter, feed *feeds.Feed) {
	atom, err := feed.ToAtom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/atom+xml")
	fmt.Fprintln(w, atom)
}

func feedHandler(fetcher func() (*feeds.Feed, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		feed, err := fetcher()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		renderFeed(w, feed)
	}
}

func main() {
	_ = godotenv.Load()

	http.HandleFunc("/puroland-info", feedHandler(GetPurolandInfo))
	http.HandleFunc("/puroland-news", feedHandler(GetPurolandNews))
	http.HandleFunc("/character-show", feedHandler(GetCharacterShow))
	http.HandleFunc("/sanrio-event", feedHandler(GetSanrioEvent))
	http.HandleFunc("/kittychan-info", feedHandler(GetKittychanInfo))
	http.HandleFunc("/sanrio-events-calendar", feedHandler(GetSanrioEventsCalendar))
	http.HandleFunc("/seibuen-event", feedHandler(GetSeibuenEvent))

	port := os.Getenv("PORT")
	if port == "" {
		port = "13000"
	}
	http.ListenAndServe(":"+port, nil)
}
