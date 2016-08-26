package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	accessPrivate      = "PRIVATE"
	accessForFollowers = "FOR_FOLLOWERS"
	accessPublic       = "PUBLIC"

	statusOffline         = "OFFLINE"
	statusOnlineAvailable = "ONLINE"
	statusUnvailable      = "UNAVAILABLE"

	repeatEveryDay            = "EVERYDAY"
	repeatEveryWeekAndCustom  = "EVERY_WEEK"
	repeatEveryMonthAndCustom = "EVERY_MONTH"

	statusSubgoalDone    = "DONE"
	statusSubgoalNotDone = "NOT_DONE"

	collectionUser  = "users"
	collectionGoal  = "goals"
	collectionHabit = "habits"
	collectionTask  = "tasks"
)

const (
	sunday = iota + 1
	monday
	tuesday
	wednesday
	thursday
	friday
	saturday
)

func validateAccessibility(access string) error {
	switch access {
	case accessPrivate, accessForFollowers, accessPublic:
		return nil
	default:
		return newError(codeBadData, fmt.Sprintf("accessibility %s is undefined (only %s, %s, and %s allowed)",
			access, accessPrivate, accessPublic, accessForFollowers))
	}
}

func validateRepeat(repeat string) error {
	switch repeat {
	case repeatEveryDay, repeatEveryMonthAndCustom, repeatEveryWeekAndCustom:
		return nil
	default:
		return newError(codeBadData, fmt.Sprintf("repeat field \"%s\" is undefined (only %s, %s, and %s allowed)",
			repeat, repeatEveryDay, repeatEveryMonthAndCustom,
			repeatEveryWeekAndCustom))
	}
}

func validateStatus(status string) error {
	switch status {
	case statusOffline, statusOnlineAvailable, statusUnvailable:
		return nil
	default:
		return newError(codeBadData, fmt.Sprintf("status %s is undefined(only %s, %s, and %s)",
			status, statusOffline, statusOnlineAvailable,
			statusUnvailable))
	}
}

func validateStatusSubgoal(status string) error {
	switch status {
	case statusSubgoalDone, statusSubgoalNotDone:
		return nil
	default:
		return newError(codeBadData, fmt.Sprintf(
			"status %s is not allowed(only %s and %s)",
			status, statusSubgoalDone, statusSubgoalNotDone))
	}
}

// Reminder represents a reminder for tasks
type Reminder struct {
	At       time.Time     `bson:"remindAt" json:"remindAt"`
	Status   string        `bson:"status" json:"status" `
	Duration time.Duration `bson:"duration" json:"duration"`
}

// Validate implements Validator's Validate
func (r *Reminder) Validate() error {
	err := validateStatusSubgoal(r.Status)
	if err != nil {
		return err
	}
	return nil
}

// RepeatReminder represents a reminder for habit
type RepeatReminder struct {
	Status            string        `bson:"status" json:"status"`
	Cycle             string        `bson:"cycle" json:"cycle"`
	DaysInWeekOrMonth map[int]bool  `bson:"days" json:"repeat_on"`
	TimeInDay         time.Duration `bson:"remindAt" json:"remindAt"`
	Duration          time.Duration `bson:"duration" json:"duration"`
}

// Validate implements Validator's Validate
func (r *RepeatReminder) Validate() error {
	err := validateStatusSubgoal(r.Status)
	if err != nil {
		return err
	}
	if err = validateRepeat(r.Cycle); err != nil {
		return err
	}
	return nil
}

// User represents a user
type User struct {
	ID                   bson.ObjectId   `bson:"_id,omitempty" json:"-"`
	Username             string          `bson:"username" json:"username"`
	Email                string          `bson:"email" json:"email"`
	Firstname            string          `bson:"firstName" json:"firstName"`
	Lastname             string          `bson:"lastName" json:"lastName"`
	Status               string          `bson:"status" json:"status"`
	DefaultAccessibility string          `bson:"defaultAccess" json:"defaultAccess"`
	CreatedAt            time.Time       `bson:"createdAt" json:"createdAt"`
	PictureURL           *string         `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty"`
	Followers            map[string]bool `bson:"followers,omitempty" json:"followers,omitempty"`
	Password             *string         `bson:"password" json:"password,omitempty"`
	Salt                 *string         `bson:"salt" json:"-"`
}

// Generate salt creates a new random salt
func generateSalt() string {
	var salt [20]byte
	rand.Seed(time.Now().UTC().UnixNano())
	dictionary := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	for i := 0; i < 20; i++ {
		salt[i] = dictionary[rand.Intn(20)]
	}
	return string(salt[:])
}

// encryptUserPassword encrypts a password using the provided salt
// Same salt and password always result in a same encrypted password (no side
// effects )
func encryptPassword(salt string, password string) string {
	var appended = []byte(password + salt)
	h := sha256.New()
	h.Write(appended)
	return hex.EncodeToString(h.Sum(nil))
}

func (u *User) encryptPassword() error {
	if u.Password == nil || len(*u.Password) == 0 {
		return newError(codeBadData, "the password not provided")
	}
	if len(*u.Password) < 6 {
		return newError(codeBadData, "the password is invalid: must be longer than or equal to 6 characters")
	}
	u.Salt = new(string)
	*u.Salt = string(generateSalt())
	*u.Password = encryptPassword(*u.Salt, *u.Password)
	return nil
}

// Validate implements the Validator interface{}
func (u *User) Validate() error {
	var err error
	switch {
	case len(u.Email) == 0:
		err = newError(codeBadData, "email is not provided")
	case len(u.Firstname) == 0:
		err = newError(codeBadData, "first name is not provided")
	case len(u.Lastname) == 0:
		err = newError(codeBadData, "last name is not provided")
	}
	if err != nil {
		return err
	}
	// encrypt user's password
	if err := u.encryptPassword(); err != nil {
		return err
	}
	// set the default accessibility
	if len(u.DefaultAccessibility) == 0 {
		u.DefaultAccessibility = accessPrivate
	}
	if err = validateAccessibility(u.DefaultAccessibility); err != nil {
		return err
	}

	u.CreatedAt = time.Now()
	return nil
}

// Cname implements Item's Cname
func (u User) Cname() string {
	return collectionUser
}

// KeySet implements Item's KeySet
func (u User) KeySet() bson.M {
	return bson.M{
		"username": u.Username,
	}
}

// SubGoal represents a common data for all goals
type SubGoal struct {
	Name        string    `bson:"name" json:"name"`
	Description *string   `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
}

// Validate implements Validator's Validate on the SubGoal Object
func (sg *SubGoal) Validate() error {
	sg.CreatedAt = time.Now()
	switch {
	case len(sg.Name) == 0:
		return newError(codeBadData, "attribute name is not provided")
	}
	return nil
}

// Goal represents a user's goal
type Goal struct {
	ID            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	*SubGoal      `bson:",inline"`
	User          string  `bson:"user" json:"user"`
	PictureURL    *string `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty"`
	Status        string  `bson:"status" json:"status"`
	Accessibility string  `bson:"accessibility" json:"accessibility,omitempty"`
}

// Validate implements Validator's Validate
func (g *Goal) Validate() error {
	if err := g.SubGoal.Validate(); err != nil {
		return err
	}
	switch {
	case len(g.Accessibility) == 0:
		return newError(codeBadData, "attribute accessibility is not provided")
	default:
		if err := validateAccessibility(g.Accessibility); err != nil {
			return err
		}
	}
	return nil
}

// Cname implements Item's Cname
func (g Goal) Cname() string {
	return collectionGoal
}

// KeySet implements Item's Cname
func (g Goal) KeySet() bson.M {
	return bson.M{
		"_id":  g.ID,
		"user": g.User,
	}
}

// Habit represents a goal's habit
type Habit struct {
	ID              bson.ObjectId `bson:"_id,omitempty" json:"id"`
	*SubGoal        `bson:",inline"`
	User            string `bson:"user" json:"-"`
	Goal            string `bson:"goal" json:"-"`
	*RepeatReminder `bson:"reminder,omitempty" json:"reminder,omitempty"`
}

// Cname implements Item's Cname
func (h Habit) Cname() string {
	return collectionHabit
}

// KeySet implements Item's Cname
func (h Habit) KeySet() bson.M {
	return bson.M{
		"_id":  h.ID,
		"user": h.User,
		"goal": h.Goal,
	}
}

// Task represents a goal's task
type Task struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	*SubGoal  `bson:",inline"`
	User      string `bson:"user" json:"-"`
	Goal      string `bson:"goal" json:"-"`
	*Reminder `bson:"reminder,omitempty" json:"reminder,omitempty"`
}

// Cname implements Item's Cname
func (t Task) Cname() string {
	return collectionTask
}

// KeySet implements Item's KeySet
func (t Task) KeySet() bson.M {
	return bson.M{
		"_id":  t.ID,
		"user": t.User,
		"goal": t.Goal,
	}
}
