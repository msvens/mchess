package model

import (
	"strings"
	"time"
)

// Date represents a date that can be parsed from multiple formats.
// schack.se API returns dates in various formats: "2025-01-22", "2025-01-22T10:00:00", etc.
type Date struct {
	time.Time
}

// UnmarshalJSON handles flexible date parsing
func (d *Date) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "" || s == "null" {
		return nil
	}

	// Try common formats
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02T15:04:05",       // ISO without timezone
		"2006-01-02",                // Date only
	}

	var err error
	for _, format := range formats {
		d.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return err
}

// MarshalJSON outputs the date in RFC3339 format
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Time.Format(time.RFC3339) + `"`), nil
}
