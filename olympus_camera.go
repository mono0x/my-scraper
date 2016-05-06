package main

import "github.com/gorilla/feeds"

const (
	OlympusCameraTitle  = "オリンパスカメラ"
	OlympusCameraUserId = "FotoPus"
)

func GetOlympusCamera() (*feeds.Feed, error) {
	posts, err := GetPostsFromFacebook(OlympusCameraUserId)
	if err != nil {
		return nil, err
	}
	return GetOlympusCameraFromPosts(posts)
}

func GetOlympusCameraFromPosts(posts *FacebookPosts) (*feeds.Feed, error) {
	return RenderFacebookFeed(posts, OlympusCameraUserId, OlympusCameraTitle)
}
