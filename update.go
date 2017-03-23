package calendarservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

type rawDate struct {
	DateTime string `json:"dateTime"`
	Date     string `json:"date"`
}

type Event struct {
	gorm.Model
	CalendarID  string
	GoogleID    string `gorm:"not null;unique"`
	Summary     string
	Description string
	Start       time.Time
	End         time.Time
}

type Calendar struct {
	ID string
}

func (event Event) IsPast() bool {
	return event.End.Sub(time.Now()) < 0
}

func log(message string) {
	fmt.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), message)
}

var Db *gorm.DB
var err error

func InitDb() {
	if Db == nil {
		Db, err = gorm.Open("sqlite3", "store.db")
		if err != nil {
			println("Database error.")
			os.Exit(1)
		}
		Db.AutoMigrate(&Event{})
	}
}

func (calendar Calendar) StartUpdate() {
	InitDb()
	for {
		rawData, err := getRawCalendarData(calendar.ID)
		if err != nil {
			log(fmt.Sprintf("Failed to fetch '%s' data: %s\n", calendar.ID, err))
		} else {
			formattedEvents := getFormattedEvents(rawData)
			addedEvents := 0
			updatedEvents := 0
			for _, event := range formattedEvents {
				event.CalendarID = calendar.ID
				existing := Event{}
				Db.Where("google_id = ?", event.GoogleID).First(&existing)
				if existing.ID == 0 {
					Db.Create(&event)
					addedEvents++
				} else {
					Db.Model(&existing).Updates(event)
					updatedEvents++
				}
			}
			log(fmt.Sprintf("%s: %d events updated, %d events created.", calendar.ID, updatedEvents, addedEvents))
		}
		time.Sleep(time.Minute)
	}
}

func getRawCalendarData(calendarID string) (rawCalendarData, error) {
	timeMin := strings.Replace(time.Now().AddDate(0, 0, -14).Format(time.RFC3339), "+", "-", 1)
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events?key=%s&singleEvents=true&timeMin=%s", calendarID, os.Getenv("GOOGLE_API_KEY"), timeMin)
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

func getFormattedEvents(data rawCalendarData) []Event {
	rawEvents := data.Items
	events := make([]Event, len(rawEvents))
	for i, rawEvent := range rawEvents {
		start, _ := parseRawDate(rawEvent.Start)
		end, _ := parseRawDate(rawEvent.End)
		events[i] = Event{GoogleID: rawEvent.ID, Summary: rawEvent.Summary, Description: rawEvent.Description, Start: start, End: end}
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Unix() < events[j].Start.Unix()
	})
	return events
}

func parseRawDate(date rawDate) (time.Time, error) {
	if date.Date != "" {
		return time.Parse("2006-01-02", date.Date)
	}
	return time.Parse("2006-01-02T15:04:05-07:00", date.DateTime)
}
