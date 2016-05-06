package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/braintree/manners"
	"github.com/gorilla/feeds"
	"github.com/joho/godotenv"
	"github.com/lestrrat/go-server-starter/listener"
)

func renderFeed(w http.ResponseWriter, feed *feeds.Feed) {
	if err := feed.WriteAtom(w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/atom+xml")
}

func feedHandler(fetcher func() (*feeds.Feed, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		feed, err := fetcher()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		renderFeed(w, feed)
	}
}

func main() {
	_ = godotenv.Load()

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
	mux.HandleFunc("/character-show", feedHandler(GetCharacterShow))
	mux.HandleFunc("/gotouchi-chara-calendar", feedHandler(GetGotouchiCharaCalendar))
	mux.HandleFunc("/kittychan-info", feedHandler(GetKittychanInfo))
	mux.HandleFunc("/memoirs-of-shibasaki-saki", feedHandler(GetMemoirsOfShibasakiSaki))
	mux.HandleFunc("/mucchan-musao", feedHandler(GetMucchanMusao))
	mux.HandleFunc("/olympus-camera", feedHandler(GetOlympusCamera))
	mux.HandleFunc("/puroland-info", feedHandler(GetPurolandInfo))
	mux.HandleFunc("/sanrio-alerts", feedHandler(GetSanrioAlerts))
	mux.HandleFunc("/sanrio-event", feedHandler(GetSanrioEvent))
	mux.HandleFunc("/sanrio-events-calendar", feedHandler(GetSanrioEventsCalendar))
	mux.HandleFunc("/seibuen-event", feedHandler(GetSeibuenEvent))

	manners.Serve(l, mux)
}
