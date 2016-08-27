package data

import (
	"fmt"
	"time"
)

// SubGoal represents a common data for all goals
type SubGoal struct {
	Name        string    `bson:"name" json:"name"`
	Description *string   `bson:"description,omitempty" json:"description,omitempty"`
	Status      string    `bson:"status" json:"status"`
	LastUpdated time.Time `bson:"createdAt" json:"createdAt"`
}

const (
	statusSubgoalDone    = "DONE"
	statusSubgoalNotDone = "NOT_DONE"
)

func validateStatusSubgoal(status string) error {
	switch status {
	case statusSubgoalDone, statusSubgoalNotDone:
		return nil
	default:
		return newBadData(fmt.Sprintf("status '%s' is not allowed(only %s and %s)",
			status, statusSubgoalDone, statusSubgoalNotDone))
	}
}

// Validate implements Validator's Validate on the SubGoal Object
func (sg *SubGoal) Validate() error {
	sg.LastUpdated = time.Now()
	err := validateStatusSubgoal(sg.Status)
	if err != nil {
		return err
	}
	switch {
	case len(sg.Name) == 0:
		return newBadData("attribute name is not provided")
	}
	return nil
}
