package server

// DefaultConfig is server default configuration
var DefaultConfig = Config{
	Domain: "127.0.0.1",
	Port:   5088,
	DBURL:  "localhost",
	DBName: "donit",
}

// Config represents a server configuration structure
type Config struct {
	// DomainName is the domain of this server without the scheme
	Domain string
	Port   int
	// The database url as defined in
	// https://godoc.org/gopkg.in/mgo.v2#Dial
	DBURL string
	// The database name
	DBName string
}
