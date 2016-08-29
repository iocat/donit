package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/iocat/donit/server"
)

var (
	port  = flag.Int("port", 5088, "the port the server listens on")
	dbURL = flag.String("db", "localhost", "the address of the database ([mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options])")
)

func main() {
	flag.Parse()
	conf := &server.Config{
		Domain: server.DefaultConfig.Domain,
		Port:   *port,
		DBURL:  *dbURL,
		DBName: server.DefaultConfig.DBName,
	}
	s, err := server.New(conf)
	if err != nil {
		fmt.Printf("%s: %s\n", os.Args[0], err)
		return
	}
	s.Start()
}
