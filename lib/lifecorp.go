package scraper

const (
	LifeCorpUserId = "lifecorp428"
)

func NewLifeCorpFacebookSource() *FacebookSource {
	return NewFacebookSource(LifeCorpUserId)
}
