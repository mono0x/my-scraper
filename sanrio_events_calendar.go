package main

const (
	SanrioEventsCalendarId  = "qsqrk2emvnnvu45debac9dugr8@group.calendar.google.com"
	SanrioEventsCalendarUrl = "https://calendar.google.com/calendar/embed?src=qsqrk2emvnnvu45debac9dugr8@group.calendar.google.com"
)

func NewSanrioEventsCalendarGoogleCalendarSource() *GoogleCalendarSource {
	return NewGoogleCalendarSource(SanrioEventsCalendarId, SanrioEventsCalendarUrl)
}
