package lalapiroomevent

import (
	"context"
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
	mux.HandleFunc("/topics/lalapiroom/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/www.lifecorp.jp/topics/lalapiroom/index.html")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client())
	source.baseURL = server.URL

	feed, err := source.Scrape(context.Background(), url.Values{})
	if err != nil {
		t.Fatal(err)
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 4, len(feed.Items))
	assert.Equal(t, "http://www.lifecorp.jp/topics/lalapiroom/#:~:text=2024%E5%B9%B44%E6%9C%8814%E6%97%A5%EF%BC%88%E6%97%A5%EF%BC%89,-%E7%A6%8F%E5%B4%8E%E5%BA%97%EF%BC%88%E5%85%B5%E5%BA%AB%EF%BC%89", feed.Items[0].Link.Href)
	assert.Equal(t, "2024年4月14日（日） 福崎店（兵庫）", feed.Items[0].Title)
	assert.WithinDuration(t, time.Date(2024, time.April, 14, 0, 0, 0, 0, loc), feed.Items[0].Created, 0)
}
