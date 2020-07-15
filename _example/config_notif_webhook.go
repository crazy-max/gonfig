package main

import (
	"time"
)

// NotifWebhook holds webhook notification configuration details
type NotifWebhook struct {
	Endpoint string            `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	Method   string            `yaml:"method,omitempty" json:"method,omitempty"`
	Headers  map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Timeout  *time.Duration    `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

// GetDefaults gets the default values
func (s *NotifWebhook) GetDefaults() *NotifWebhook {
	n := &NotifWebhook{}
	n.SetDefaults()
	return n
}

// SetDefaults sets the default values
func (s *NotifWebhook) SetDefaults() {
	s.Method = "GET"
	s.Timeout = NewDuration(10 * time.Second)
}
