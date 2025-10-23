package utils

import (
	"log"
	"time"
)

var LOC *time.Location

func LoadTimezone() {
	var err error
	LOC, err = time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Printf("Failed to load timezone: %v", err)
	}
}
