package fukokulifeevent

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
	mux.HandleFunc("/event/index.php", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/act.fukoku-life.co.jp/event/index.php")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client())
	source.baseURL = server.URL

	feed, err := source.Scrape()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 6, len(feed.Items))
	assert.Equal(t, "https://act.fukoku-life.co.jp/event/detail.php?id=415", feed.Items[0].Link.Href)
	assert.Equal(t, "フコク生命 ピラティス体験イベント（in京浜）", feed.Items[0].Title)
	assert.Contains(t, feed.Items[0].Description, "2016年10月12日（水）")
	assert.Contains(t, feed.Items[0].Description, "川崎市川崎区貝塚1-1-3 川崎フコク生命ビル４Ｆ 特設会場")
	assert.Equal(t, "https://act.fukoku-life.co.jp/event/detail.php?id=435", feed.Items[5].Link.Href)
}
