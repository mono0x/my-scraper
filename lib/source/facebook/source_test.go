package facebook

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	source := NewSource(http.DefaultClient, "ACCESS_TOKEN", "user")
	assert.Equal(t, http.DefaultClient, source.httpClient)
	assert.Equal(t, "ACCESS_TOKEN", source.accessToken)
	assert.Equal(t, "user", source.userID)
	assert.Equal(t, baseURL, source.baseURL)
}

func TestScrape(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v2.6/mucchan.musao/posts", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "ACCESS_TOKEN", query.Get("access_token"))
		assert.Equal(t, "created_time,from,link,message,picture", query.Get("fields"))
		http.ServeFile(w, r, "testdata/graph.facebook.com/v2.6/mucchan.musao/posts")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client(), "ACCESS_TOKEN", "mucchan.musao")
	source.baseURL = server.URL

	feed, err := source.Scrape()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "高尾山公認キャラ\u3000ムッちゃん", feed.Title)
	assert.Equal(t, 25, len(feed.Items))
	assert.Equal(t, "★☆Next ムサさび〜ず☆★", feed.Items[0].Title)
	assert.Equal(t, "高尾山公認キャラ\u3000ムッちゃん", feed.Items[0].Author.Name)
	assert.Equal(t, "https://www.facebook.com/mucchan.musao/videos/1124654310918067/", feed.Items[0].Link.Href)
	assert.Equal(t, "https://www.facebook.com/mucchan.musao/posts/1123833604333471/", feed.Items[1].Link.Href)
}
