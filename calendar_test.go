package calendarservice

import "testing"
import "fmt"
import "os"

var calendar = Calendar{ID: os.Getenv("CALENDAR_ID"), Name: "test"}

func TestEventsByWeeks(t *testing.T) {
	calendar.Update()
	eventsByWeeks := calendar.EventsByWeeks()
	for i, weeklyEvents := range eventsByWeeks {
		fmt.Printf("WEEK %d\n", i)
		for _, event := range weeklyEvents {
			fmt.Printf("\t%s: %s\n", event.Start, event.Summary)
		}
	}
}

func TestEventByWeekdayName(t *testing.T) {
	eventsByWeeks := calendar.EventsByWeeks()
	for i, weeklyEvents := range eventsByWeeks {
		fmt.Printf("WEEK %d\n", i)
		fmt.Printf("\tFRIDAY:\t\t%s\n", weeklyEvents.EventByWeekdayName("friday").Summary)
		fmt.Printf("\tSATURDAY:\t%s\n", weeklyEvents.EventByWeekdayName("saturday").Summary)
	}
}

func TestCalendarUpdate(t *testing.T) {
}
