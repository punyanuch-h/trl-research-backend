package cron

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// NewCronWithTimezone creates a cron instance set to Asia/Bangkok
func NewCronWithTimezone() *cron.Cron {
	jakartaTime, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("⚠️ Could not load Asia/Bangkok location, falling back to local: %v", err)
		jakartaTime = time.Local
	}
	return cron.New(cron.WithLocation(jakartaTime))
}
