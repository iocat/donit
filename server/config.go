package server

// DefaultConfig is server default configuration
var DefaultConfig = Config{
	Domain: "api.donit.xyz",
}

// Config represents a server configuration structure
type Config struct {
	// DomainName is the domain of this server without the scheme
	Domain string
}
