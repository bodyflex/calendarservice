package calendarservice

import (
	"fmt"
	"time"
)

var week = time.Duration(int(time.Hour) * 24 * 7)

var weekdayIndexes = map[string]int{
	"Monday":    0,
	"Tuesday":   1,
	"Wednesday": 2,
	"Thursday":  3,
	"Friday":    4,
	"Saturday":  5,
	"Sunday":    6,
}

type WeeklyEvents map[string]Event

func (weeklyEvents WeeklyEvents) EventsByWeekdays(weekdays ...string) []Event {
	events := make([]Event, 0)
	for _, weekday := range weekdays {
		event := weeklyEvents[weekday]
		events = append(events, event)
	}
	return events
}

type Calendar struct {
	Name string
	ID   string
}

func (this Calendar) FilledCalendar(startWeekOffset int, endWeekOffset int, weekdays []string) [][]Event {
	eventData := this.EventsByWeeks(fmt.Sprintf("%d week", startWeekOffset), fmt.Sprintf("%d week", endWeekOffset))
	weeklyEvents := make([][]Event, 0)
	for i := 0; i < endWeekOffset-startWeekOffset; i++ {
		events := make([]Event, 0)
		for _, weekday := range weekdays {
			var event Event
			found := true
			if i < len(eventData) {
				event, found = eventData[i][weekday]
			} else {
				found = false
			}
			if !found {
				start := time.Now().Add(time.Duration((startWeekOffset + i) * int(week))).Truncate(week).Add(time.Duration(weekdayIndexes[weekday] * 24 * int(time.Hour)))
				event = Event{Start: start}
			}
			events = append(events, event)
		}
		weeklyEvents = append(weeklyEvents, events)
	}
	return weeklyEvents
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

func (this Calendar) EventsByWeeks(startInterval string, endInterval string) []WeeklyEvents {
	initDb()

	allEvents := []Event{}
	db.
		Where("start >= date_trunc('week', now()) + ?::interval AND start < date_trunc('week', now()) + ?::interval AND calendar_name = ?", startInterval, endInterval, this.Name).
		Order("start asc").
		Find(&allEvents)

	weeklyEventsList := make([]WeeklyEvents, 0)
	var weeklyEvents WeeklyEvents
	var previousWeekNumber int
	for _, event := range allEvents {
		_, weekNumber := event.Start.ISOWeek()
		if weekNumber != previousWeekNumber {
			if weeklyEvents != nil {
				weeklyEventsList = append(weeklyEventsList, weeklyEvents)
			}
			weeklyEvents = WeeklyEvents{}
			previousWeekNumber = weekNumber
		}
		weeklyEvents[event.Start.Weekday().String()] = event
	}
	weeklyEventsList = append(weeklyEventsList, weeklyEvents)
	return weeklyEventsList
}
