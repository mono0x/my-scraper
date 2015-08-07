package main

import (
	"bufio"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetSanrioEventFromDocument(t *testing.T) {
	f, err := os.Open("data/www.sanrio.co.jp/event/search/index.html")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	feed, err := GetSanrioEventFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(feed.Items), 20)
}

