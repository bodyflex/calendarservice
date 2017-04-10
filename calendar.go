package calendarservice

import (
	"fmt"
	"strings"
	"time"
)

type Calendar struct {
	Name string
	ID   string
}

type WeeklyEvents []Event

func (weeklyEvents WeeklyEvents) EventByWeekdayName(weekday string) Event {
	for _, event := range weeklyEvents {
		if strings.ToLower(event.Start.Weekday().String()) == strings.ToLower(weekday) {
			return event
		}
	}
	return Event{}
}

func (this Calendar) StartUpdate() {
	for {
		this.Update()
		time.Sleep(time.Minute)
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
			events = make(WeeklyEvents, 0)
			previousWeekNumber = weekNumber
		}
		events = append(events, event)
	}
	weeklyEvents = append(weeklyEvents, events)
	return weeklyEvents
}
