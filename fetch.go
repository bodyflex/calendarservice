package calendarservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type rawCalendarData struct {
	ETag  string     `json:"etag"`
	Items []rawEvent `json:"items"`
}

type rawEvent struct {
	ID          string  `json:"id"`
	Kind        string  `json:"kind"`
	Status      string  `json:"status"`
	Summary     string  `json:"summary"`
	Description string  `json:"description"`
	Start       rawDate `json:"start"`
	End         rawDate `json:"end"`
}

func (this rawEvent) getFormattedEvent() Event {
	start, _ := this.Start.parse()
	end, _ := this.End.parse()
	return Event{GoogleID: this.ID, Summary: this.Summary, Description: this.Description, Start: start, End: end}
}

type rawDate struct {
	DateTime string `json:"dateTime"`
	Date     string `json:"date"`
}

func (this rawDate) parse() (time.Time, error) {
	if this.Date != "" {
		return time.Parse("2006-01-02", this.Date)
	}
	return time.Parse("2006-01-02T15:04:05-07:00", this.DateTime)
}

func fetchRawCalendarData(calendarID string) (rawCalendarData, error) {
	timeMin := strings.Replace(time.Now().AddDate(0, 0, -14).Format(time.RFC3339), "+", "-", 1)
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events?key=%s&singleEvents=true&showDeleted=true&timeMin=%s", calendarID, os.Getenv("GOOGLE_API_KEY"), timeMin)

	response, err := http.Get(url)
	if err != nil {
		return rawCalendarData{}, err
	}

	defer response.Body.Close()
	calendar := rawCalendarData{}
	if err = json.NewDecoder(response.Body).Decode(&calendar); err != nil {
		return rawCalendarData{}, err
	}
	return calendar, nil
}

func log(message string) {
	fmt.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), message)
}
