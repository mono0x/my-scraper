package main

import (
	"bufio"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetPurolandNewsFromDocument(t *testing.T) {
	f, err := os.Open("data/puroland.jp/index.html")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	feed, err := GetPurolandNewsFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, len(feed.Items))
}

func TestGetPurolandInfoFromDocument(t *testing.T) {
	f, err := os.Open("data/puroland.jp/index.html")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	feed, err := GetPurolandInfoFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 50, len(feed.Items))
}
