package server

// DefaultConfig is server default configuration
var DefaultConfig = Config{
	Domain: "localhost",
	Port:   5088,
}

// Config represents a server configuration structure
type Config struct {
	// DomainName is the domain of this server without the scheme
	Domain string
	Port   int
}
