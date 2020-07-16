package config

import "time"

// ServerFTP holds ftp server configuration
type ServerFTP struct {
	Host               string         `yaml:"host,omitempty" json:"host,omitempty" validate:"required"`
	Port               int            `yaml:"port,omitempty" json:"port,omitempty" validate:"required,min=1"`
	Username           string         `yaml:"username,omitempty" json:"username,omitempty"`
	Password           string         `yaml:"password,omitempty" json:"password,omitempty"`
	Sources            []string       `yaml:"sources,omitempty" json:"sources,omitempty"`
	Timeout            *time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	DisableEPSV        *bool          `yaml:"disableEPSV,omitempty" json:"disableEPSV,omitempty"`
	TLS                *bool          `yaml:"tls,omitempty" json:"tls,omitempty"`
	InsecureSkipVerify *bool          `yaml:"insecureSkipVerify,omitempty" json:"insecureSkipVerify,omitempty"`
	LogTrace           *bool          `yaml:"logTrace,omitempty" json:"logTrace,omitempty"`
}

// GetDefaults gets the default values
func (s *ServerFTP) GetDefaults() *ServerFTP {
	n := &ServerFTP{}
	n.SetDefaults()
	return n
}

// SetDefaults sets the default values
func (s *ServerFTP) SetDefaults() {
	s.Port = 21
	s.Sources = []string{}
	s.Timeout = NewDuration(5 * time.Second)
	s.DisableEPSV = NewFalse()
	s.TLS = NewFalse()
	s.InsecureSkipVerify = NewFalse()
	s.LogTrace = NewFalse()
}
