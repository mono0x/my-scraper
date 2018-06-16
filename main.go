package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/lestrrat-go/server-starter/listener"
	"github.com/mono0x/my-scraper/lib/server"
	"github.com/pkg/errors"
)

func run() error {
	listeners, err := listener.ListenAll()
	if err != nil && err != listener.ErrNoListeningTarget {
		return errors.WithStack(err)
	}

	var l net.Listener
	if len(listeners) > 0 {
		l = listeners[0]
	} else {
		l, err = net.Listen("tcp", ":8080")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	mux := server.NewServeMux()
	server := http.Server{Handler: mux}

	go func() {
		if err := server.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)

	for {
		s := <-signalChan
		if s == syscall.SIGTERM {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		}
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	_ = godotenv.Load()

	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
