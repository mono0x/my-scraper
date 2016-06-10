package scraper

const (
	OlympusCameraUserId = "FotoPus"
)

func NewOlympusCameraFacebookSource() *FacebookSource {
	return NewFacebookSource(OlympusCameraUserId)
}
