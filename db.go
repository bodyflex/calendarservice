package calendarservice

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

func initDb() {
	if db == nil {
		db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			panic(err)
		}
		db.AutoMigrate(&Event{})
	}
}
