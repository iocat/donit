package data

import (
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// User represents a user
type User struct {
	Username             string    `bson:"username" json:"username"`
	Email                string    `bson:"email" json:"email"`
	Firstname            string    `bson:"firstName" json:"firstName"`
	Lastname             string    `bson:"lastName" json:"lastName"`
	Status               string    `bson:"status" json:"status"`
	DefaultAccessibility string    `bson:"defaultAccess" json:"defaultAccess"`
	LastUpdated          time.Time `bson:"lastUpdated" json:"lastUpdated"`
	PictureURL           *string   `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty"`
	Password             *string   `bson:"password" json:"password,omitempty"`
	Salt                 *string   `bson:"salt" json:"-"`
}

const (
	statusOffline         = "OFFLINE"
	statusOnlineAvailable = "ONLINE"
	statusUnvailable      = "UNAVAILABLE"
)

func validateStatus(status string) error {
	switch status {
	case statusOffline, statusOnlineAvailable, statusUnvailable:
		return nil
	default:
		return newBadData(fmt.Sprintf("status %s is undefined(only %s, %s, and %s)",
			status, statusOffline, statusOnlineAvailable,
			statusUnvailable))
	}
}

// generateSalt creates a new random salt
func generateSalt() string {
	var salt [20]byte
	rand.Seed(time.Now().UTC().UnixNano())
	dictionary := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	for i := 0; i < 20; i++ {
		salt[i] = dictionary[rand.Intn(20)]
	}
	return string(salt[:])
}

// encryptPassword encrypts the password
func (u *User) encryptPassword() error {
	if u.Password == nil || len(*u.Password) == 0 {
		return newBadData("the password not provided")
	}
	if len(*u.Password) < 6 {
		return newBadData("the password is invalid: must be longer than or equal to 6 characters")
	}
	u.Salt = new(string)
	*u.Salt = string(generateSalt())
	*u.Password = encryptPassword(*u.Salt, *u.Password)
	return nil
}

// Validate implements Validator's Validate
func (u *User) Validate() error {
	var err error
	switch {
	case len(u.Email) == 0:
		err = newBadData("email is not provided")
	case len(u.Firstname) == 0:
		err = newBadData("first name is not provided")
	case len(u.Lastname) == 0:
		err = newBadData("last name is not provided")
	}
	if err != nil {
		return err
	}
	if err := u.encryptPassword(); err != nil {
		return err
	}
	if err = validateAccessibility(u.DefaultAccessibility); err != nil {
		return err
	}
	if err = validateStatus(u.Status); err != nil {
		return err
	}
	u.LastUpdated = time.Now()
	return nil
}

// Cname implements Item's Cname
func (u *User) cname() string {
	return CollectionUser
}

// KeySet implements Item's keys
func (u *User) keys() bson.M {
	return bson.M{
		"username": u.Username,
	}
}

// subKeys get the sub keys and associated collection
func (u *User) subCols() (keys bson.M, cname []string) {
	k := bson.M{
		"user": u.Username,
	}
	return k, []string{CollectionFollower, CollectionGoal, CollectionHabit, CollectionTask}
}

// SetKeys implements Item's SetKeys
func (u *User) SetKeys(k []string) error {
	if err := writeKey(k, &u.Username); err != nil {
		return err
	}
	return nil
}
