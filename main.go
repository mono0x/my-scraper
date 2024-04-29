package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mono0x/my-scraper/scraper"
	"github.com/mono0x/my-scraper/scraper/source/googlecalendar"
	"github.com/mono0x/my-scraper/scraper/source/kittychaninfo"
	"github.com/mono0x/my-scraper/scraper/source/lalapiroomevent"
	"github.com/mono0x/my-scraper/scraper/source/yuyakekoyakenews"
	"github.com/mono0x/my-scraper/server"
	"golang.org/x/sync/errgroup"
)

func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	sources := []scraper.Source{
		googlecalendar.NewSource(client),
		kittychaninfo.NewSource(client),
		lalapiroomevent.NewSource(client),
		yuyakekoyakenews.NewSource(client),
	}

	handler, err := server.NewHandler(sources)
	if err != nil {
		return err
	}
	s := http.Server{Handler: handler}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	eg := errgroup.Group{}
	eg.Go(func() error {
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("%w", err)
		}
		return nil
	})
	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	})
	return eg.Wait()
}

func main() {
	log.SetFlags(log.Lshortfile)

	_ = godotenv.Load()

	if err := run(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
