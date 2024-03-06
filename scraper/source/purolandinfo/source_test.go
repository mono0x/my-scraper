package purolandinfo

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	mux.HandleFunc("/api/live/get_information/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/www.puroland.jp/api/live/get_information/index.json")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client())
	source.baseURL = server.URL

	feed, err := source.Scrape(url.Values{})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 5, len(feed.Items))
}
