package users

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/errors"
	"github.com/iocat/donit/internal/donitdoc/internal/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CollectionName for secondary indexing purpose
var CollectionName = utils.User.MGOName()

// col is collection getter function for user
var col func(*mgo.Database) *mgo.Collection

func init() {
	col = utils.User.Collection
}

// CreateDoc creates a new user, the provided password must not be encrypted.
// CreateDoc also validates the user's data before doing as insertion
// package user uses consistent encryption algorithm to encrypts the password
func CreateDoc(l log.Logger, context string, db *mgo.Database, user *User, password string) error {
	l.Log("ctx", context, "op", "users.CreateDoc", "user", user.Username)
	// Encrypt the password
	encrypted, salt, err := randomSaltEncryption(password)
	if err != nil {
		l.Log("ctx", context, "user", user, "result", fmt.Errorf("validation error: encryption: %s", err))
		return errors.NewValidate(err.Error())
	}
	// Create a stored user
	var suser = StoredUser{
		User:     *user,
		Password: encrypted,
		Salt:     salt,
	}
	// Evaluate user's data consistency
	if err := utils.Validate(suser); err != nil {
		l.Log("ctx", context, "user", user, "result", fmt.Errorf("validation error: %s", err))
		return errors.NewValidate(err.Error())
	}
	// Insert a stored user
	if err := col(db).Insert(suser); err != nil {
		if mgo.IsDup(err) {
			l.Log("ctx", context, "result", "DUPLICATED")
			return errors.NewDuplicated("user", user.Username)
		}
		l.Log("ctx", context, "result", fmt.Sprintf("error: %s", err))
		return err
	}
	l.Log("ctx", context, "result", "SUCCESS")
	return nil
}

// UpdateDoc replaces a user's data (not included password and salt) and validates before replacement
func UpdateDoc(l log.Logger, context string, db *mgo.Database, user *User) error {
	l.Log("ctx", context, "op", "users.UpdateDoc", "user", user.Username)
	// Validate a user data
	if err := utils.Validate(user); err != nil {
		l.Log("ctx", context, "result", fmt.Errorf("validation error: %s", err))
		return errors.NewValidate(err.Error())
	}
	// Upsert a user without replacing the password and salt
	err := col(db).Update(bson.M{
		"username": user.Username,
	}, bson.M{
		"$set": user,
	})
	if err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", context, "result", "NOT_FOUND")
			return errors.NewNotFound("user", user.Username)
		}
		l.Log("ctx", context, "result", fmt.Sprintf("error: %s", err))
		return err
	}
	l.Log("ctx", context, "result", "SUCCESS")
	return nil
}

// DeleteDoc deletes a user
func DeleteDoc(l log.Logger, context string, db *mgo.Database, username string) error {
	l.Log("ctx", context, "op", "users.DeleteDoc", "user", username)
	if err := col(db).Remove(bson.M{
		"username": username,
	}); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", context, "result", "NOT_FOUND")
			return errors.NewNotFound("user", username)
		}
		l.Log("ctx", context, "result", fmt.Sprintf("error: %s", err))
		return err
	}
	l.Log("ctx", context, "result", "SUCCESS")
	return nil
}

// ReadDoc reads a user data
func ReadDoc(l log.Logger, context string, db *mgo.Database, username string) (*User, error) {
	l.Log("ctx", context, "op", "users.ReadDoc", "user", username)
	var user StoredUser
	if err := col(db).Find(bson.M{
		"username": username,
	}).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", context, "result", "NOT_FOUND")
			return nil, errors.NewNotFound("user", username)
		}
		l.Log("ctx", context, "result", fmt.Errorf("error: %s", err))
		return nil, err
	}
	l.Log("ctx", context, "result", "SUCCESS")
	return &user.User, nil
}

// ValidatePassword validates the user's password
// TODO: inefficient: allocate a whole user to store just password + salt
func ValidatePassword(l log.Logger, context string, db *mgo.Database, username, password string) (bool, error) {
	l.Log("ctx", context, "op", "users.ValidatePassword", "user", username)
	var user StoredUser
	if err := col(db).Find(bson.M{
		"username": username,
	}).Select(bson.M{
		"password": 1,
		"salt":     1,
	}).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", context, "result", "NOT_FOUND")
			return false, errors.NewNotFound("user", username)
		}
		l.Log("ctx", context, "result", fmt.Sprintf("error: %s", err))
		return false, err
	}
	if user.Password != encryptPassword(user.Salt, password) {
		l.Log("ctx", context, "result", "SUCCESS, WRONG PASSWORD")
		return false, nil
	}
	l.Log("ctx", context, "result", "SUCCESS")
	return true, nil
}

// ChangePassword changes the user's password. Caller needs to provide an old password for
// authentication before changing the user's password
func ChangePassword(l log.Logger, context string, db *mgo.Database, username, old, new string) error {
	// Evaluate old password
	ok, err := ValidatePassword(l, context, db, username, old)
	if err != nil {
		l.Log("ctx", context, "op", "users.ChangePassword", "user", username, "result", "validate old password: FAILED")
		return err
	}
	// Old password is incorrect
	if !ok {
		l.Log("ctx", context, "result", "validate old password: FAILED")
		return errors.NewValidate("wrong old password")
	}
	// Encrypt the new password
	encrypted, salt, err := randomSaltEncryption(new)
	if err != nil {
		l.Log("ctx", context, "result", fmt.Sprintf("validation error: encryption %s", err))
		return errors.NewValidate(err.Error())
	}
	// Change the password
	if err := col(db).Update(bson.M{
		"username": username,
	}, bson.M{
		"$set": bson.M{
			"password": encrypted,
			"salt":     salt,
		},
	}); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", context, "result", "NOT FOUND")
			return errors.NewNotFound("user", username)
		}
		l.Log("ctx", context, "result", fmt.Sprintf("update password on document: %s", err))
		return err
	}
	l.Log("ctx", context, "result", "SUCCESS")
	return nil
}
