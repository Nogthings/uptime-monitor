package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Exclude from JSON responses
	CreatedAt time.Time `json:"created_at"`
}

// HashPassword hashes the plaintext password and sets it to the User's Password field.
func (u *User) HashPassword(plaintextPassword string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword compares the provided plaintext password with the stored hashed password.
func (u *User) CheckPassword(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plaintextPassword))
	if err != nil {
		switch {
		case err == bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
