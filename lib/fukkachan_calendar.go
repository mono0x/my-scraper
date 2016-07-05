package scraper

const (
	FukkachanCalendarId  = "fukkachan.com@gmail.com"
	FukkachanCalendarUrl = "https://calendar.google.com/calendar/embed?src=fukkachan.com@gmail.com"
)

func NewFukkachanCalendarGoogleCalendarSource() *GoogleCalendarSource {
	return NewGoogleCalendarSource(FukkachanCalendarId, FukkachanCalendarUrl)
}
