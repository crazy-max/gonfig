package config

import (
	"time"
)

// NewFalse returns a false bool pointer
func NewFalse() *bool {
	b := false
	return &b
}

// NewTrue returns a true bool pointer
func NewTrue() *bool {
	b := true
	return &b
}

// NewDuration returns a duration pointer
func NewDuration(duration time.Duration) *time.Duration {
	return &duration
}
