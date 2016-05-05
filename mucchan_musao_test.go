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

	assert.Equal(t, 25, len(feed.Items))
	assert.Equal(t, "★☆Next ムサさび〜ず☆★", feed.Items[0].Title)
}
