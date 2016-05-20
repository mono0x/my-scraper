package main

const (
	GotouchiCharaCalendarId  = "gnr0r3kevuuv3j0q6q25gj4hks@group.calendar.google.com"
	GotouchiCharaCalendarUrl = "https://calendar.google.com/calendar/embed?src=gnr0r3kevuuv3j0q6q25gj4hks@group.calendar.google.com"
)

func NewGotouchiCharaCalendarGoogleCalendarSource() *GoogleCalendarSource {
	return NewGoogleCalendarSource(GotouchiCharaCalendarId, GotouchiCharaCalendarUrl)
}
