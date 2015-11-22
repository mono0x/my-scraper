package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetPurolandInfoFromDocument(t *testing.T) {
	f, err := os.Open("data/www.puroland.jp/api/live/get_information/index.json")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	feed, err := GetPurolandInfoFromReader(f)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 5, len(feed.Items))
}
