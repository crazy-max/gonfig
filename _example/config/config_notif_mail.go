package config

// NotifMail holds mail notification configuration details
type NotifMail struct {
	Host               string `yaml:"host,omitempty" json:"host,omitempty"`
	Port               int    `yaml:"port,omitempty" json:"port,omitempty"`
	SSL                *bool  `yaml:"ssl,omitempty" json:"ssl,omitempty"`
	InsecureSkipVerify *bool  `yaml:"insecureSkipVerify,omitempty" json:"insecureSkipVerify,omitempty"`
	Username           string `yaml:"username,omitempty" json:"username,omitempty"`
	UsernameFile       string `yaml:"usernameFile,omitempty" json:"usernameFile,omitempty"`
	Password           string `yaml:"password,omitempty" json:"password,omitempty"`
	PasswordFile       string `yaml:"passwordFile,omitempty" json:"passwordFile,omitempty"`
	From               string `yaml:"from,omitempty" json:"from,omitempty"`
	To                 string `yaml:"to,omitempty" json:"to,omitempty"`
}

// GetDefaults gets the default values
func (s *NotifMail) GetDefaults() *NotifMail {
	n := &NotifMail{}
	n.SetDefaults()
	return n
}

// SetDefaults sets the default values
func (s *NotifMail) SetDefaults() {
	s.Host = "localhost"
	s.Port = 25
	s.SSL = NewFalse()
	s.InsecureSkipVerify = NewTrue()
}
