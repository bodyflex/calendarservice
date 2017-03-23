package calendarservice

import "testing"
import "fmt"

func TestEventsByWeeks(t *testing.T) {
	eventsByWeeks := EventsByWeeks()
	for i, weeklyEvents := range eventsByWeeks {
		fmt.Printf("WEEK %d\n", i)
		for _, event := range weeklyEvents {
			fmt.Printf("\t%s: %s\n", event.Start, event.Summary)
		}
	}
}

func TestEventByWeekdayName(t *testing.T) {
	eventsByWeeks := EventsByWeeks()
	for i, weeklyEvents := range eventsByWeeks {
		fmt.Printf("WEEK %d\n", i)
		fmt.Printf("\tFRIDAY:\t\t%s\n", weeklyEvents.EventByWeekdayName("friday").Summary)
		fmt.Printf("\tSATURDAY:\t%s\n", weeklyEvents.EventByWeekdayName("saturday").Summary)
	}
}
