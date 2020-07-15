package main

// Db holds data necessary for database configuration
type Db struct {
	Path string `yaml:"path,omitempty" json:"path,omitempty"`
}

// GetDefaults gets the default values
func (s *Db) GetDefaults() *Db {
	n := &Db{}
	n.SetDefaults()
	return n
}

// SetDefaults sets the default values
func (s *Db) SetDefaults() {
	s.Path = "bbolt.db"
}
