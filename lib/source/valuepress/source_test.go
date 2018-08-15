package valuepress

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
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "サンリオ", query.Get("q"))
		http.ServeFile(w, r, "testdata/www.value-press.com/search")
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

	assert.Equal(t, 20, len(feed.Items))
	assert.Equal(t, "ハローキティも出演！世界初！お皿に盛りつけられたコスメを使って楽しむ参加型メイクアップレッスンショー『ミワンダフルのメイクアップレストラン』が開催", feed.Items[0].Title)
	assert.Equal(t, "https://www.value-press.com/pressrelease/171658", feed.Items[0].Link.Href)
	assert.Equal(t, "メイクスマイルプロジェクト所属のメイクスマイルアーティスト・ミワンダフルがお届けする、参加型メイクアップレッスンショー『ミワンダフルのメイクアップレストラン』が2016年11月22(火)、 23日(水・祝)に原宿クエストホールにてオープンします。", feed.Items[0].Description)
	assert.Equal(t, "メイクスマイルプロジェクト", feed.Items[0].Author.Name)
	assert.WithinDuration(t, time.Date(2016, 10, 17, 10, 0, 0, 0, loc), feed.Items[0].Created, 0)
}
