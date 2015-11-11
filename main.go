package main

import (
	"fmt"
	"github.com/braintree/manners"
	"github.com/gorilla/feeds"
	"github.com/lestrrat/go-server-starter/listener"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM)
	go func() {
		for {
			s := <-signalChan
			if s == syscall.SIGTERM {
				manners.Close()
			}
		}
	}()

	listeners, err := listener.ListenAll()
	if err != nil {
		log.Fatal(err)
	}

	var l net.Listener
	if len(listeners) > 0 {
		l = listeners[0]
	} else {
		l, err = net.Listen("tcp", ":13000")
		if err != nil {
			log.Fatal(err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/puroland-info", feedHandler(GetPurolandInfo))
	mux.HandleFunc("/puroland-news", feedHandler(GetPurolandNews))
	mux.HandleFunc("/character-show", feedHandler(GetCharacterShow))
	mux.HandleFunc("/sanrio-event", feedHandler(GetSanrioEvent))
	mux.HandleFunc("/kittychan-info", feedHandler(GetKittychanInfo))
	mux.HandleFunc("/sanrio-events-calendar", feedHandler(GetSanrioEventsCalendar))
	mux.HandleFunc("/seibuen-event", feedHandler(GetSeibuenEvent))

	manners.Serve(l, mux)
}
