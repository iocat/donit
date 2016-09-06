package utils

import (
	"context"
	"errors"
)

// GetLogID gets the log's uuid from the context
func GetLogID(c context.Context) (string, error) {
	if id, ok := c.Value("log").(string); ok {
		return id, nil
	}
	return "", errors.New("logging uuid (marked as \"log\") is not provided in the context")
}

// MustGetLogID panics when error occurs
func MustGetLogID(c context.Context) string {
	id, err := GetLogID(c)
	if err != nil {
		panic(err)
	}
	return id
}
