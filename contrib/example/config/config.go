package config

// Config holds configuration details
type Config struct {
	Timezone string    `yaml:"timezone,omitempty" json:"timezone,omitempty"`
	LogLevel string    `yaml:"logLevel,omitempty" json:"logLevel,omitempty"`
	LogJSON  bool      `yaml:"logJSON,omitempty" json:"logJSON,omitempty"`
	Db       *Db       `yaml:"db,omitempty" json:"db,omitempty"`
	Server   *Server   `yaml:"server,omitempty" json:"server,omitempty"`
	Download *Download `yaml:"download,omitempty" json:"download,omitempty" validate:"required"`
	Notif    *Notif    `yaml:"notif,omitempty" json:"notif,omitempty"`
}
