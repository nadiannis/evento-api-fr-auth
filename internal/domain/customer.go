package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Password password `json:"-"`
	Balance  float64  `json:"balance"`
}

type password []byte

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	*p = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(*p, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
