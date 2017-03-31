package calendarservice

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	CalendarName string `gorm:"not null"`
	CalendarID   string `gorm:"not null"`
	GoogleID     string `gorm:"not null"`
	Summary      string
	Description  string
	Start        time.Time `gorm:"not null"`
	End          time.Time `gorm:"not null"`
}

func (this Event) IsPast() bool {
	return this.End.Sub(time.Now()) < 0
}
