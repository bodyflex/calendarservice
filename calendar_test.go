package calendarservice

import "testing"

import "os"

var calendar = Calendar{ID: os.Getenv("CALENDAR_ID"), Name: "test"}

func TestIterator(t *testing.T) {
	calendar.Update()
	calendar.FilledCalendar(0, 4, []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"})
}

func TestEventsByWeeks(t *testing.T) {
	t.Skip()
	// calendar.Update()
	// eventsByWeeks := calendar.EventsByWeeks("0 week", "1 week")
	// for i, weeklyEvents := range eventsByWeeks {
	// 	fmt.Printf("WEEK %d\n", i)
	// 	for _, event := range weeklyEvents {
	// 		fmt.Printf("\t%s: %s\n", event.Start, event.Summary)
	// 	}
	// }
}

func TestEventByWeekdayName(t *testing.T) {
	t.Skip()
	// eventsByWeeks := calendar.EventsByWeeks("0 week", "1 week")
	// for i, weeklyEvents := range eventsByWeeks {
	// 	fmt.Printf("WEEK %d\n", i)
	// }
}

func TestCalendarUpdate(t *testing.T) {
}
