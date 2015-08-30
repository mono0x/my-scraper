package main

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetKittychanInfoFromDocument(t *testing.T) {
	f, err := os.Open("data/www.kittychan.info/information.html")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	feed, err := GetKittychanInfoFromReader(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 100, len(feed.Items))
}
