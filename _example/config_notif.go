package main

// Notif holds data necessary for notification configuration
type Notif struct {
	Mail    *NotifMail    `yaml:"mail,omitempty" json:"mail,omitempty"`
	Webhook *NotifWebhook `yaml:"webhook,omitempty" json:"webhook,omitempty"`
}
