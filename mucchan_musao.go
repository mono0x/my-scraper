package main

const (
	MucchanMusaoUserId = "mucchan.musao"
)

func NewMucchanMusaoFacebookSource() *FacebookSource {
	return NewFacebookSource(MucchanMusaoUserId)
}
