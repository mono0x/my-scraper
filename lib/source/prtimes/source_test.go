package prtimes

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
	mux.HandleFunc("/main/action.php", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/prtimes.jp/main/action.php")
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

	assert.Equal(t, 40, len(feed.Items))
	assert.Equal(t, "今年のテーマは“ギフト”！サンリオキャラクターが贈る冬のイベント「ピューロウィンターギフト」開催決定！", feed.Items[0].Title)
	assert.Equal(t, "https://prtimes.jp/main/html/rd/p/000000085.000007643.html", feed.Items[0].Link.Href)
	assert.Equal(t, "株式会社サンリオエンターテイメント", feed.Items[0].Author.Name)
	assert.WithinDuration(t, time.Date(2016, 9, 29, 14, 3, 3, 0, loc), feed.Items[0].Created, 0)
}
