package scraper

const (
	MucchanMusaoUserId = "mucchan.musao"
)

func NewMucchanMusaoFacebookSource() *FacebookSource {
	return NewFacebookSource(MucchanMusaoUserId)
}
