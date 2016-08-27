package data

import "gopkg.in/mgo.v2/bson"

// Goal represents a user's goal
type Goal struct {
	ID            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	SubGoal       `bson:"subGoal,inline"`
	User          string  `bson:"user" json:"user"`
	PictureURL    *string `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty"`
	Accessibility string  `bson:"accessibility" json:"accessibility,omitempty"`
}

// Validate implements Validator's Validate
func (g *Goal) Validate() error {
	if err := g.SubGoal.Validate(); err != nil {
		return err
	}
	switch {
	case len(g.Accessibility) == 0:
		return newBadData("attribute accessibility is not provided")
	default:
		if err := validateAccessibility(g.Accessibility); err != nil {
			return err
		}
	}
	return nil
}

// cname implements Item's cname
func (g *Goal) cname() string {
	return CollectionGoal
}

// keys implements Item's keys
func (g *Goal) keys() bson.M {
	return bson.M{
		"_id":  g.ID,
		"user": g.User,
	}
}

func (g *Goal) subCols() (keys bson.M, cname []string) {
	k := bson.M{
		"goal": g.ID,
	}
	return k, []string{CollectionHabit, CollectionTask, CollectionComment}
}

// SetKeys implements Item's SetKeys
func (g *Goal) SetKeys(k []string) error {
	if err := writeKey(k, &g.User, &g.ID); err != nil {
		return err
	}
	return nil
}
