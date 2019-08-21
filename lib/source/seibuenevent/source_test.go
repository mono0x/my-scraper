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
		http.ServeFile(w, r, "testdata/www.seibu-leisure.co.jp/event/12410/index.html")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client())
	source.baseURL = server.URL

	feed, err := source.Scrape()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(feed.Items))
	assert.Equal(t, "7月、8月メルヘンタウンイベント", feed.Items[0].Title)
	assert.Equal(t, "https://www.seibu-leisure.co.jp/event/19387/index.html", feed.Items[0].Link.Href)
	assert.Equal(t, "★メルヘンタウンのお姉さん★", feed.Items[1].Title)
	assert.Equal(t, "https://www.seibu-leisure.co.jp/event/melhen_sister/index.html", feed.Items[1].Link.Href)
}
