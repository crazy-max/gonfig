package config

// Server represents a server configuration
type Server struct {
	FTP *ServerFTP `yaml:"ftp,omitempty" json:"ftp,omitempty"`
}

// GetDefaults gets the default values
func (s *Server) GetDefaults() *Server {
	return nil
}

// SetDefaults sets the default values
func (s *Server) SetDefaults() {
	// noop
}
