package facebook

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/mono0x/my-scraper/lib"
	"github.com/stretchr/testify/assert"
)

var _ scraper.Source = (*Source)(nil)

func TestSource(t *testing.T) {
	file, err := os.Open("testdata/graph.facebook.com/v2.6/mucchan.musao/posts")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var posts Posts
	if err := json.NewDecoder(file).Decode(&posts); err != nil {
		t.Fatal(err)
	}

	source := NewSource("mucchan.musao")
	feed, err := source.Render(&posts)
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
