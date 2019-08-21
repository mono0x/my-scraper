package harmonylandinfo

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
	mux.HandleFunc("/welcome.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/www.harmonyland.jp/welcome.html")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client())
	source.baseURL = server.URL

	feed, err := source.Scrape()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 13, len(feed.Items))
	assert.Equal(t, "2017年はシナモン15周年！                     ハーモニーランドはシナモンイベントがいっぱい♪【6/2～7/14】", feed.Items[0].Title)
	assert.Equal(t, "http://www.harmonyland.jp/event/rain/index.html", feed.Items[0].Link.Href)
}
