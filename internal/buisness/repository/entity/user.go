package entity

import (
	"github.com/AgeroFlynn/crud/internal/buisness/core/dto"
	"time"
)

import "github.com/lib/pq"

// User represents an individual user.
type User struct {
	ID           string         `pg:"user_id,pk,type:uuid"`
	Name         string         `pg:"name"`
	Email        string         `pg:"email"`
	Roles        pq.StringArray `pg:"roles"`
	PasswordHash []byte         `pg:"password_hash"`
	DateCreated  time.Time      `pg:"date_created"`
	DateUpdated  time.Time      `pg:"date_updated"`
}

func (u *User) ToDTOUser() *dto.User {
	return &dto.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Roles:        u.Roles,
		PasswordHash: u.PasswordHash,
		DateCreated:  u.DateCreated,
		DateUpdated:  u.DateUpdated,
	}
}

func FromDTOUser(user *dto.User) *User {
	return &User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Roles:        user.Roles,
		PasswordHash: user.PasswordHash,
		DateCreated:  user.DateCreated,
		DateUpdated:  user.DateUpdated,
	}
}

func ToDTOUserSlice(users *[]User) *[]dto.User {
	var dtoUsers []dto.User

	for _, user := range *users {
		dtoUsers = append(dtoUsers, *user.ToDTOUser())
	}
	return &dtoUsers
}
