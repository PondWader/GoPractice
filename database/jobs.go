package database

import (
	"time"

	"gorm.io/gorm"
)

func runDbJobs(db *gorm.DB) {
	db.Where("time_cached < ?", time.Now().Add(-time.Hour)).Delete(&SessionCache{})

	time.Sleep(time.Minute * 5)
	runDbJobs(db)
}
