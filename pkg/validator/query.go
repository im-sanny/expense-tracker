// package validator provides query parameter validation utilities
package validator

import (
	"strconv"
	"strings"
	"time"
)

const (
	// default pagination value
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
	DateLayout   = "2006-01-01"
)

// ParseInt parses and validates an integer with optional min/max bounds
func ParseInt(val string, defaultValue, minBound, maxBound int) (int, error) {
	if val == "" {
		return defaultValue, nil
	}

	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	if minBound > 0 && n < minBound {
		return defaultValue, nil
	}

	if maxBound > 0 && n > maxBound {
		return defaultValue, nil
	}

	return n, nil
}

// ParseDate parses a date string, returns zero time if empty
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	return time.Parse(DateLayout, dateStr)
}

// ParseSearch parses and trims search query
func ParseSearch(q string) string {
	return strings.TrimSpace(q)
}
