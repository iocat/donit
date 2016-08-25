package main

import (
	"fmt"
	"os"

	"github.com/iocat/donit/server"
)

func main() {
	s, err := server.New(nil)
	if err != nil {
		fmt.Printf("%s:%s", os.Args[0], err)
	}
	s.Start()
}
