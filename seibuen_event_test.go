package main

import (
	"bufio"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetSeibuenEventFromDocument(t *testing.T) {
	f, err := os.Open("data/www.seibuen-yuuenchi.jp/event/index.html")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	feed, err := GetSeibuenEventFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 6, len(feed.Items))
}
