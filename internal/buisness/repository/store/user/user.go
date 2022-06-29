// Package user contains user related CRUD functionality.
package user

import (
	"context"
	"fmt"
	"github.com/AgeroFlynn/crud/internal/buisness/core/dto"
	"github.com/AgeroFlynn/crud/internal/buisness/repository/entity"
	"github.com/AgeroFlynn/crud/internal/buisness/sys/auth"
	"github.com/AgeroFlynn/crud/internal/buisness/sys/validate"
	"github.com/AgeroFlynn/crud/internal/foundation/database"
	"github.com/go-pg/pg/v10"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Store manages the set of API's for user access.
type Store struct {
	log *zap.SugaredLogger
	db  *pg.DB
}

// NewStore constructs a user store for api access.
func NewStore(log *zap.SugaredLogger, db *pg.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database.
func (s Store) Create(ctx context.Context, nu dto.NewUser, now time.Time) (dto.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.User{}, fmt.Errorf("generating password hash: %w", err)
	}

	usr := entity.User{
		ID:           validate.GenerateID(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		DateCreated:  now,
		DateUpdated:  now,
	}

	_, err = s.db.Model(&usr).Insert()
	if err != nil {
		return dto.User{}, fmt.Errorf("inserting user: %w", err)
	}

	return *usr.ToDTOUser(), nil
}

// Update replaces a user document in the database.
func (s Store) Update(ctx context.Context, claims auth.Claims, userID string, uu dto.UpdateUser, now time.Time) error {
	usr, err := s.FindByID(ctx, claims, userID)
	if err != nil {
		return fmt.Errorf("updating user userID[%s]: %w", userID, err)
	}

	if uu.Name != nil {
		usr.Name = *uu.Name
	}
	if uu.Email != nil {
		usr.Email = *uu.Email
	}
	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}
	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}
		usr.PasswordHash = pw
	}
	usr.DateUpdated = now

	if _, err := s.db.Model(&usr).Where("user_id = ?", userID).Update(); err != nil {
		return fmt.Errorf("updating userID[%s]: %w", userID, err)
	}

	return nil
}

// Delete removes a user from the database.
func (s Store) Delete(ctx context.Context, claims auth.Claims, userID string) error {
	// If you are not an admin and looking to delete someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
		return database.ErrForbidden
	}

	var usr entity.User
	if _, err := s.db.Model(&usr).Where("user_id = ?", userID).Delete(); err != nil {
		return fmt.Errorf("deleting userID[%s]: %w", userID, err)
	}

	return nil
}

// FindAll retrieves a list of existing users from the database.
func (s Store) FindAll(ctx context.Context) ([]dto.User, error) {

	var users []entity.User
	if err := s.db.Model(&users).Select(); err != nil {
		if err == pg.ErrNoRows {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("selecting users: %w", err)
	}

	return *entity.ToDTOUserSlice(&users), nil
}

// FindByID gets the specified user from the database.
func (s Store) FindByID(ctx context.Context, claims auth.Claims, userID string) (dto.User, error) {
	if err := validate.CheckID(userID); err != nil {
		return dto.User{}, database.ErrInvalidID
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
		return dto.User{}, database.ErrForbidden
	}

	var usr entity.User
	if err := s.db.Model(&usr).Where("user_id = ?", userID).Limit(1).Select(); err != nil {
		if err == pg.ErrNoRows {
			return dto.User{}, database.ErrNotFound
		}
		return dto.User{}, fmt.Errorf("selecting userID[%q]: %w", userID, err)
	}

	return *usr.ToDTOUser(), nil
}

// FindByEmail gets the specified user from the database by email.
func (s Store) FindByEmail(ctx context.Context, claims auth.Claims, email string) (dto.User, error) {

	var usr entity.User
	if err := s.db.Model(&usr).Where("email = ?", email).Limit(1).Select(); err != nil {
		if err == pg.ErrNoRows {
			return dto.User{}, database.ErrNotFound
		}
		return dto.User{}, fmt.Errorf("selecting email[%q]: %w", email, err)
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != usr.ID {
		return dto.User{}, database.ErrForbidden
	}

	return *usr.ToDTOUser(), nil
}

// Authenticate finds a user by their email and verifies their password. On
// success, it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (s Store) Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {

	var usr entity.User
	if err := s.db.Model(&usr).Where("email = ?", email).Limit(1).Select(); err != nil {
		if err == pg.ErrNoRows {
			return auth.Claims{}, database.ErrNotFound
		}
		return auth.Claims{}, fmt.Errorf("selecting user[%q]: %w", email, err)
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function - so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, database.ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "service project",
			Subject:   usr.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: usr.Roles,
	}

	return claims, nil
}
