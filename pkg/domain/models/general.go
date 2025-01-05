package models

import (
	"time"
)

const (
	LayoutTimestamp       = time.RFC3339 // format "2006-01-02T15:04:05Z07:00"
	MaxAttempts           = 11
	MinRangeSleepDuration = 100 * time.Millisecond // min range time wait offset
	MaxRangeSleepDuration = 500 * time.Millisecond // max range time wait offset
	SleepOffset           = 50 * time.Millisecond  // offset
	SaveOffset            = 10
	MaxRowsFromDB         = 999
	MaxTimeoutContext     = 600 * time.Second
	TimeDriftForExpire    = 600 // 10 minutes
	MaxTimeForLocks       = 30 * time.Second
	TimeoutRequest        = 5 * time.Minute
)

// 1970-01-01T00:00:00
var TimeDefault = time.Unix(0, 0) // time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC)

const (
	CredentialCreateContextKey   = "createcredential"
	CredentialExchangeContextKey = "exchangecredential"
	ActionGoogleKey              = "actiongoogle"
	CommandTypeCreate            = "create"
	CommandTypeUpdate            = "update"
	CommandTypeDelete            = "delete"
)
