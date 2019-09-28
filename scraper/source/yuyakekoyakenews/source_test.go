package yuyakekoyakenews

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
	mux.HandleFunc("/news/index.php", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/yuyakekoyake.jp/news/index.php")
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

	assert.Equal(t, 15, len(feed.Items))
	assert.Equal(t, "http://yuyakekoyake.jp/news/news_detail.php?id=news59f16132092f0", feed.Items[0].Link.Href)
	assert.Equal(t, "イルミネーションの実施につきまして", feed.Items[0].Title)
	assert.WithinDuration(t, time.Date(2017, 10, 26, 0, 0, 0, 0, loc), feed.Items[0].Created, 0)
}
