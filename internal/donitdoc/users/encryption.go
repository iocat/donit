package users

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"

	stdliberr "errors"
)

// generateSalt creates a new random salt for password encryption
func generateSalt() string {
	var salt [20]byte
	rand.Seed(time.Now().UTC().UnixNano())
	dictionary := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	for i := 0; i < 20; i++ {
		salt[i] = dictionary[rand.Intn(20)]
	}
	return string(salt[:])
}

// encryptPassword encrypts a password using the provided salt
// Same salt and password always result in a same encrypted password (no side
// effects )
// The password encryption's algorithm is SHA256
func encryptPassword(salt string, password string) string {
	var appended = []byte(password + salt)
	h := sha256.New()
	h.Write(appended)
	return hex.EncodeToString(h.Sum(nil))
}

// randomSaltEncryption encrypts the password with a random salt and returns
// the encrypted password and the salt, in that order.
// randomSalt encryption uses encryptPassword to generate the encrypted password
// password must be long enough, or an error shall be thrown
func randomSaltEncryption(password string) (string, string, error) {
	if len(password) < 6 {
		return "", "", stdliberr.New("invalid password, password must be longer than 6")
	}
	salt := generateSalt()
	encrypted := encryptPassword(salt, password)
	return encrypted, salt, nil
}
