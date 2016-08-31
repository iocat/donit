package utils

import "fmt"

const (
	// User is the collection index for User
	User = iota
	// Goal is the collection index for Goal
	Goal
	// Habit is the collection index for Habit
	Habit
	// Task is the collection index for Task
	Task
	// Comment is the collection index for Comment
	Comment
)

// mgoCollectionName contains the names of the collection that store
// the corresponding data in mongodb
var mgoCollectionNames = []string{
	User:    "users",
	Goal:    "goals",
	Habit:   "habits",
	Task:    "tasks",
	Comment: "comments",
}

// GetMgoCollection returns the collection name corresponding to the index value
// Index is the index name coded in this package
func GetMgoCollection(code int) (string, error) {
	if 0 <= code && code < len(mgoCollectionNames) {
		if len(mgoCollectionNames[code]) == 0 {
			return "", fmt.Errorf("collection name corresponding to index %d is not found", code)
		}
		return mgoCollectionNames[code], nil
	}
	return "", fmt.Errorf("collection index is out of range, got index %d, expected range is [%d,%d]",
		code, 0, len(mgoCollectionNames)-1)
}
