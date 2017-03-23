package calendarservice

import "strings"

type WeeklyEvents []Event

func (weeklyEvents WeeklyEvents) EventByWeekdayName(weekday string) Event {
	for _, event := range weeklyEvents {
		if strings.ToLower(event.Start.Weekday().String()) == strings.ToLower(weekday) {
			return event
		}
	}
	return Event{}
}

func (calendar Calendar) EventsByWeeks() []WeeklyEvents {
	InitDb()

	allEvents := []Event{}
	Db.
		Where("start > date('now', '-7 day', 'weekday 1')").
		Where("calendar_id = ?", calendar.ID).
		Order("start asc").
		Find(&allEvents)

	weeklyEvents := make([]WeeklyEvents, 0)
	var events WeeklyEvents
	var previousWeekNumber int
	for _, event := range allEvents {
		_, weekNumber := event.Start.ISOWeek()
		if weekNumber != previousWeekNumber {
			if events != nil {
				weeklyEvents = append(weeklyEvents, events)
			}
			events = make(WeeklyEvents, 0)
			previousWeekNumber = weekNumber
		}
		events = append(events, event)
	}
	return weeklyEvents
}
