package seibuenevent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	source := NewSource(http.DefaultClient)
	assert.Equal(t, http.DefaultClient, source.httpClient)
	assert.Equal(t, baseURL, source.baseURL)
}

func TestScrape(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "e1", query.Get("category"))
		http.ServeFile(w, r, "testdata/www.seibu-leisure.co.jp/event/index.html")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client())
	source.baseURL = server.URL

	feed, err := source.Scrape()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(feed.Items))
}
