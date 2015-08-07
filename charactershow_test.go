package main

import (
	"bufio"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetCharacterShowFromDocument(t *testing.T) {
	f, err := os.Open("data/charactershow.jp/index.html")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	feed, err := GetCharacterShowFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(feed.Items), 115)
}
