package utils

import (
	"context"
	"errors"
)

// UUID represents the UUID for logging and tracing purpose
type UUID string

// GetLogID gets the log's uuid from the context
func GetLogID(c context.Context) (UUID, error) {
	if id, ok := c.Value("log").(UUID); ok {
		return id, nil
	}
	return "", errors.New("logging uuid (marked as \"log\") is not provided in the context")
}

// MustGetLogID panics when error occurs
func MustGetLogID(c context.Context) UUID {
	id, err := GetLogID(c)
	if err != nil {
		panic(err)
	}
	return id
}
