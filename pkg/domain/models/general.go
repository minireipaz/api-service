package models

import "time"

const (
	LayoutTimestamp  = "2006-01-02T15:04:05Z07:00"
	MaxAttempts      = 11
	MinSleepDuration = 100 * time.Millisecond // min time wait
	MaxSleepDuration = 500 * time.Millisecond // max time wait
	SleepOffset      = 50 * time.Millisecond  // offset
	SaveOffset       = 10
)
