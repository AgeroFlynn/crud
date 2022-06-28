//Package usergrp represents incoming data models from request
package incoming

import (
	"github.com/AgeroFlynn/crud/internal/buisness/core/dto"
	"github.com/lib/pq"
	"time"
)

// User represents an individual user.
type User struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Roles        pq.StringArray `json:"roles"`
	PasswordHash []byte         `json:"-"`
	DateCreated  time.Time      `json:"date_created"`
	DateUpdated  time.Time      `json:"date_updated"`
}

func (u *User) ToDTOUser() dto.User {
	return dto.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Roles:        u.Roles,
		PasswordHash: u.PasswordHash,
		DateCreated:  u.DateCreated,
		DateUpdated:  u.DateUpdated,
	}
}

func FromDTOUser(user dto.User) User {
	return User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Roles:        user.Roles,
		PasswordHash: user.PasswordHash,
		DateCreated:  user.DateCreated,
		DateUpdated:  user.DateUpdated,
	}
}

func ToDTOUserSlice(users []User) []dto.User {
	var dtoUsers []dto.User

	for _, user := range users {
		dtoUsers = append(dtoUsers, user.ToDTOUser())
	}
	return dtoUsers
}

func FromDTOUserSlice(users []dto.User) []User {
	var incomingUsers []User

	for _, user := range users {
		incomingUsers = append(incomingUsers, FromDTOUser(user))
	}
	return incomingUsers
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Roles           []string `json:"roles" validate:"required"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"password_confirm" validate:"eqfield=Password"`
}

func (nu *NewUser) ToDTONewUser() dto.NewUser {
	return dto.NewUser{
		Name:            nu.Name,
		Email:           nu.Email,
		Roles:           nu.Roles,
		Password:        nu.Password,
		PasswordConfirm: nu.PasswordConfirm,
	}
}

func FromDTONewUser(nu dto.NewUser) NewUser {
	return NewUser{
		Name:            nu.Name,
		Email:           nu.Email,
		Roles:           nu.Roles,
		Password:        nu.Password,
		PasswordConfirm: nu.PasswordConfirm,
	}
}

func ToDTONewUserSlice(users []NewUser) []dto.NewUser {
	var dtoNewUsers []dto.NewUser

	for _, user := range users {
		dtoNewUsers = append(dtoNewUsers, user.ToDTONewUser())
	}
	return dtoNewUsers
}

// UpdateUser defines what information may be provided to modify an existing
// User. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields - so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types, but we make exceptions around
// marshalling/unmarshalling.
type UpdateUser struct {
	Name            *string  `json:"name"`
	Email           *string  `json:"email" validate:"omitempty,email"`
	Roles           []string `json:"roles"`
	Password        *string  `json:"password"`
	PasswordConfirm *string  `json:"password_confirm" validate:"omitempty,eqfield=Password"`
}

func (uu *UpdateUser) ToDTOUpdateUser() dto.UpdateUser {
	return dto.UpdateUser{
		Name:            uu.Name,
		Email:           uu.Email,
		Roles:           uu.Roles,
		Password:        uu.Password,
		PasswordConfirm: uu.PasswordConfirm,
	}
}

func FromDTOUpdateUser(nu dto.UpdateUser) UpdateUser {
	return UpdateUser{
		Name:            nu.Name,
		Email:           nu.Email,
		Roles:           nu.Roles,
		Password:        nu.Password,
		PasswordConfirm: nu.PasswordConfirm,
	}
}

func ToDTOUpdateUserSlice(users []UpdateUser) []dto.UpdateUser {
	var dtoNewUsers []dto.UpdateUser

	for _, user := range users {
		dtoNewUsers = append(dtoNewUsers, user.ToDTOUpdateUser())
	}
	return dtoNewUsers
}
