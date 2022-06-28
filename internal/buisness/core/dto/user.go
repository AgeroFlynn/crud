package dto

import (
	"github.com/lib/pq"
	"time"
)

// User represents an individual user.
type User struct {
	ID           string
	Name         string
	Email        string
	Roles        pq.StringArray
	PasswordHash []byte
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string
	Email           string
	Roles           []string
	Password        string
	PasswordConfirm string
}

// UpdateUser defines what information may be provided to modify an existing
// User. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateUser struct {
	Name            *string
	Email           *string
	Roles           []string
	Password        *string
	PasswordConfirm *string
}
