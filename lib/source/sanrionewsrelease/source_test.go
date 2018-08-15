package sanrionewsrelease

import (
	"net/http"
	"net/http/httptest"
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
	mux.HandleFunc("/corporate/release/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/www.sanrio.co.jp/corporate/release/index.html")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client())
	source.baseURL = server.URL

	feed, err := source.Scrape()
	if err != nil {
		t.Fatal(err)
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 51, len(feed.Items))
	assert.Equal(t, "ぐでぐでやる気のない「ぐでたま」のイベント九州初上陸！ 夏休み企画 「ぐでたま in ふくおか」 7月21日(木)〜 福岡パルコ & sanrio vivitix 天神地下街店にて開催 (PDF)", feed.Items[0].Title)
	assert.Equal(t, "http://www.sanrio.co.jp/wp-content/uploads/2015/05/20160708-1.pdf", feed.Items[0].Link.Href)
	assert.WithinDuration(t, time.Date(2016, 7, 8, 0, 0, 0, 0, loc), feed.Items[0].Created, 0)
	assert.Equal(t, "2016年バレンタイン向けスペシャルギフト「GODIVA &ハローキティ」・「GODIVA &マイメロディ」1月6日（水）よりサンリオ限定販売", feed.Items[50].Title)
	assert.Equal(t, "http://www.sanrio.co.jp/corporate/release/y2016/d0106/", feed.Items[50].Link.Href)
	assert.WithinDuration(t, time.Date(2016, 1, 6, 0, 0, 0, 0, loc), feed.Items[50].Created, 0)
}
