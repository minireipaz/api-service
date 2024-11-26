package models

import (
	"log"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == `""` || s == `"0001-01-01T00:00:00Z"` || s == "1970-01-01 00:00:00" || s == "\"1970-01-01 00:00:00\"" {
		*ct = CustomTime{Time: time.Time{}}
		return nil
	}

	s = s[1 : len(s)-1]
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		parsedTime, err = time.Parse(time.DateTime, s)
		if err != nil {
			log.Printf("ERROR | Cannot parse time: %v", err)
			*ct = CustomTime{Time: time.Time{}}
			return nil
		}
	}

	ct.Time = parsedTime
	return nil
}
