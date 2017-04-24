package calendarservice

import (
	"fmt"
	"time"
)

type Calendar struct {
	Name string
	ID   string
}

type WeeklyEvents map[int]Event

func (weeklyEvents WeeklyEvents) EventsByWeekdays(weekdays ...int) []Event {
	events := make([]Event, 0)
	for _, weekday := range weekdays {
		event := weeklyEvents[weekday]
		events = append(events, event)
	}
	return events
}

func (this Calendar) StartUpdate() {
	this.StartUpdateWithTime(time.Minute)
}

func (this Calendar) StartUpdateWithTime(duration time.Duration) {
	for {
		this.Update()
		time.Sleep(duration)
	}
}

func (this Calendar) Update() {
	initDb()
	rawData, err := fetchRawCalendarData(this.ID)
	if err != nil {
		log(fmt.Sprintf("Failed to fetch '%s' data: %s\n", this.ID, err))
	} else {
		addedEvents := 0
		updatedEvents := 0
		deletedEvents := 0
		for _, rawEvent := range rawData.Items {
			event := rawEvent.getFormattedEvent()
			event.CalendarID = this.ID
			event.CalendarName = this.Name
			existing := Event{}
			db.
				Where("calendar_name = ? and google_id = ?", event.CalendarName, event.GoogleID).
				First(&existing)
			if rawEvent.Status == "cancelled" {
				if existing.ID != 0 {
					db.Unscoped().Delete(Event{}, "google_id = ?", rawEvent.ID)
					deletedEvents++
				}
			} else if existing.ID == 0 {
				db.Create(&event)
				addedEvents++
			} else {
				db.Model(&existing).Updates(event)
				updatedEvents++
			}
		}
		log(fmt.Sprintf("%s: %d events updated, %d events created, %d events deleted.", this.ID, updatedEvents, addedEvents, deletedEvents))
	}
}

func (this Calendar) EventsByWeeks() []WeeklyEvents {
	initDb()

	allEvents := []Event{}
	db.
		Where("start > date_trunc('week', now()) and calendar_name = ?", this.Name).
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
			events = WeeklyEvents{}
			previousWeekNumber = weekNumber
		}
		events[int(event.Start.Weekday())] = event
	}
	weeklyEvents = append(weeklyEvents, events)
	return weeklyEvents
}
