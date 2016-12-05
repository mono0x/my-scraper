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

func sourceRenderer(source scraper.Source) func(http.ResponseWriter, *http.Request) {
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

func run() error {
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
		return err
	}

	var l net.Listener
	if len(listeners) > 0 {
		l = listeners[0]
	} else {
		l, err = net.Listen("tcp", ":13000")
		if err != nil {
			return err
		}
	}

	mux := http.NewServeMux()

	entries := []struct {
		Path   string
		Source scraper.Source
	}{
		{"/character-show", scraper.NewCharacterShowSource()},
		{"/fukkachan-calendar", scraper.NewGoogleCalendarSource("fukkachan.com@gmail.com")},
		{"/fukoku-life", scraper.NewFukokuLifeEventSource()},
		{"/gotouchi-chara-calendar", scraper.NewGoogleCalendarSource("gnr0r3kevuuv3j0q6q25gj4hks@group.calendar.google.com")},
		{"/harmonyland-info", scraper.NewHarmonylandInfoSource()},
		{"/kittychan-info", scraper.NewKittychanInfoSource()},
		{"/lifecorp", scraper.NewFacebookSource("lifecorp428")},
		{"/memoirs-of-shibasaki-saki", scraper.NewMemoirsOfShibasakiSakiSource()},
		{"/mucchan-musao", scraper.NewFacebookSource("mucchan.musao")},
		{"/olympus-camera", scraper.NewFacebookSource("FotoPus")},
		{"/prtimes-sanrio", scraper.NewPRTimesSource()},
		{"/puroland-info", scraper.NewPurolandInfoSource()},
		{"/sanrio-events-calendar", scraper.NewGoogleCalendarSource("qsqrk2emvnnvu45debac9dugr8@group.calendar.google.com")},
		{"/sanrio-news-release", scraper.NewSanrioNewsReleaseSource()},
		{"/seibuen-event", scraper.NewSeibuenEventSource()},
		{"/value-press-sanrio", scraper.NewValuePressSource()},
		{"/yufuterashima-calendar", scraper.NewGoogleCalendarSource("pompomyufu@gmail.com")},
	}
	for _, entry := range entries {
		mux.HandleFunc(entry.Path, sourceRenderer(entry.Source))
	}

	manners.Serve(l, mux)
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)

	_ = godotenv.Load()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
