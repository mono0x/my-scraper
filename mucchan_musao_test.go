package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMucchanMusaoFromPosts(t *testing.T) {
	jsonData, err := ioutil.ReadFile("data/graph.facebook.com/v2.6/mucchan.musao/posts")
	if err != nil {
		t.Fatal(err)
	}

	var posts FacebookPosts
	if err := json.Unmarshal(jsonData, &posts); err != nil {
		t.Fatal(err)
	}

	feed, err := GetMucchanMusaoFromPosts(&posts)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "高尾山公認キャラ\u3000ムッちゃん", feed.Title)
	assert.Equal(t, 25, len(feed.Items))
	assert.Equal(t, "★☆Next ムサさび〜ず☆★", feed.Items[0].Title)
	assert.Equal(t, "高尾山公認キャラ\u3000ムッちゃん", feed.Items[0].Author.Name)
	assert.Equal(t, "https://www.facebook.com/mucchan.musao/videos/1124654310918067/", feed.Items[0].Link.Href)
	assert.Equal(t, "https://www.facebook.com/mucchan.musao/posts/1123830127667152/", feed.Items[1].Link.Href)
}
