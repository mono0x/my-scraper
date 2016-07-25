package scraper

const (
	YufuTerashimaCalendarId  = "pompomyufu@gmail.com"
	YufuTerashimaCalendarUrl = "https://calendar.google.com/calendar/embed?src=pompomyufu@gmail.com"
)

func NewYufuTerashimaCalendarGoogleCalendarSource() *GoogleCalendarSource {
	return NewGoogleCalendarSource(YufuTerashimaCalendarId, YufuTerashimaCalendarUrl)
}
