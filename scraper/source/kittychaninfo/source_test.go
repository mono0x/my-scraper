package kittychaninfo

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
	mux.HandleFunc("/information.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/www.kittychan.info/information.html")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client())
	source.baseURL = server.URL

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	feed, err := source.Scrape(t.Context(), url.Values{})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 100, len(feed.Items))
	assert.Equal(t, "DAISO池袋東武店「ハローキティ」グリーティングイベント（ワンショット撮影会）開催（2024年5月1日）", feed.Items[0].Title)
	assert.Equal(t, "https://www.daiso-sangyo.co.jp/info/news/29560", feed.Items[0].Link.Href)
	assert.WithinDuration(t, time.Date(2024, time.May, 1, 0, 0, 0, 0, loc), feed.Items[0].Created, 0)

	assert.Equal(t, "DAISO心斎橋店　ハローキティ　グリーティングイベント（ワンショット撮影会）開催（2024年4月19日）", feed.Items[99].Title)
	assert.Equal(t, "https://www.daiso-sangyo.co.jp/info/news/29407", feed.Items[99].Link.Href)
	assert.WithinDuration(t, time.Date(2024, time.April, 19, 0, 0, 0, 0, loc), feed.Items[99].Created, 0)
}
