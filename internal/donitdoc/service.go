package donitdoc

import (
	"github.com/go-kit/kit/log"
	"gopkg.in/mgo.v2"
)

// Service is the exported interface for donit's document store
type Service interface {
}

// service implements Service interface
type service struct {
	d *mgo.Database
	l log.Logger
}
