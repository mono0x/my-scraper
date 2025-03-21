package impresswatchcolumn

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	source := NewSource(http.DefaultClient)
	assert.Equal(t, http.DefaultClient, source.httpClient)
	assert.Equal(t, baseURL, source.baseURL)
}

func TestScrape(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/docs/column/stapa/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/k-tai.watch.impress.co.jp/docs/column/stapa/index.html")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client())
	source.baseURL = server.URL

	v := url.Values{}
	v.Set("site", "k-tai")
	v.Set("column", "stapa")
	feed, err := source.Scrape(t.Context(), v)
	if err != nil {
		t.Fatal(err)
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "スタパ齋藤の「スタパトロニクスMobile」", feed.Title)
	assert.Equal(t, "https://k-tai.watch.impress.co.jp/docs/column/stapa/", feed.Link.Href)
	assert.Equal(t, 16, len(feed.Items))
	assert.Equal(t, "https://k-tai.watch.impress.co.jp/docs/column/stapa/1585879.html", feed.Items[0].Link.Href)
	assert.Equal(t, "最新のMacBook Airをついに購入！\u3000M3搭載で持ち運べる、ジョリーグッドでアルティメットなマシンだゼ!!!", feed.Items[0].Title)
	assert.WithinDuration(t, time.Date(2024, time.April, 22, 0, 0, 0, 0, loc), feed.Items[0].Created, 0)
}
