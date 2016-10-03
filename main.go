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
	"github.com/mono0x/my-scraper/lib"
)

func renderFeed(w http.ResponseWriter, feed *feeds.Feed) {
	if err := feed.WriteAtom(w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/atom+xml")
}

func renderSource(source scraper.Source) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		feed, err := source.Scrape()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		renderFeed(w, feed)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

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

	entries := []struct {
		Path   string
		Source scraper.Source
	}{
		{"/character-show", scraper.NewCharacterShowSource()},
		{"/fukkachan-calendar", scraper.NewFukkachanCalendarGoogleCalendarSource()},
		{"/gotouchi-chara-calendar", scraper.NewGotouchiCharaCalendarGoogleCalendarSource()},
		{"/kittychan-info", scraper.NewKittychanInfoSource()},
		{"/lifecorp", scraper.NewLifeCorpFacebookSource()},
		{"/memoirs-of-shibasaki-saki", scraper.NewMemoirsOfShibasakiSakiSource()},
		{"/mucchan-musao", scraper.NewMucchanMusaoFacebookSource()},
		{"/olympus-camera", scraper.NewOlympusCameraFacebookSource()},
		{"/prtimes-sanrio", scraper.NewPRTimesSource()},
		{"/puroland-info", scraper.NewPurolandInfoSource()},
		{"/sanrio-events-calendar", scraper.NewSanrioEventsCalendarGoogleCalendarSource()},
		{"/sanrio-news-release", scraper.NewSanrioNewsReleaseSource()},
		{"/seibuen-event", scraper.NewSeibuenEventSource()},
		{"/yufuterashima-calendar", scraper.NewYufuTerashimaCalendarGoogleCalendarSource()},
	}
	for _, entry := range entries {
		mux.HandleFunc(entry.Path, renderSource(entry.Source))
	}

	manners.Serve(l, mux)
}
