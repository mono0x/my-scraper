package main

import "github.com/gorilla/feeds"

const (
	MucchanMusaoUserId = "mucchan.musao"
)

func GetMucchanMusao() (*feeds.Feed, error) {
	posts, err := GetPostsFromFacebook(MucchanMusaoUserId)
	if err != nil {
		return nil, err
	}
	return GetMucchanMusaoFromPosts(posts)
}

func GetMucchanMusaoFromPosts(posts *FacebookPosts) (*feeds.Feed, error) {
	return RenderFacebookFeed(posts, MucchanMusaoUserId)
}
