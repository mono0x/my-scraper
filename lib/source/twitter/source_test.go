package twitter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	source := NewSource(http.DefaultClient, 725638238291943424)
	assert.Equal(t, http.DefaultClient, source.httpClient)
	assert.Equal(t, int64(725638238291943424), source.userID)
	assert.Equal(t, anaconda.BaseUrl, source.baseURL)
}

func TestScrape(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/statuses/user_timeline.json", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "725638238291943424", query.Get("user_id"))
		assert.Equal(t, "100", query.Get("count"))
		http.ServeFile(w, r, "testdata/api.twitter.com/1.1/statuses/user_timeline.json")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client(), 725638238291943424)
	source.baseURL = server.URL

	feed, err := source.Scrape()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Sanrio Events (@sanrio_events)", feed.Title)
	assert.Equal(t, 20, len(feed.Items))
	assert.Equal(t, "12/24 KT 藤崎本館 (藤崎本館): https://t.co/Gh7HkC6yM2", feed.Items[0].Title)
	assert.Equal(t, "https://twitter.com/sanrio_events/status/805532146110582784", feed.Items[0].Link.Href)
	assert.WithinDuration(t, time.Date(2016, 12, 4, 22, 0, 4, 0, time.UTC), feed.Items[0].Created, 0)
}
