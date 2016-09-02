package utils

import (
	"context"
)

// UUID represents the UUID for logging and tracing purpose
type UUID string

// GetLogID gets the log's uuid from the context
func GetLogID(c context.Context) UUID {
	return c.Value("log").(UUID)
}
