// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package user

import (
	"context"
	"fmt"
	"github.com/AgeroFlynn/crud/internal/buisness/core/dto"
	"github.com/AgeroFlynn/crud/internal/buisness/repository/user"
	"github.com/AgeroFlynn/crud/internal/buisness/sys/auth"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
	"time"
)

// Core manages the set of API's for user access.
type Core struct {
	log  *zap.SugaredLogger
	user user.Store
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, db *pg.DB) Core {
	return Core{
		log:  log,
		user: user.NewStore(log, db),
	}
}

// Create inserts a new user into the database.
func (c Core) Create(ctx context.Context, nu dto.NewUser, now time.Time) (dto.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	usr, err := c.user.Create(ctx, nu, now)
	if err != nil {
		return dto.User{}, fmt.Errorf("create: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}

// Update replaces a user document in the database.
func (c Core) Update(ctx context.Context, claims auth.Claims, userID string, uu dto.UpdateUser, now time.Time) error {

	// PERFORM PRE BUSINESS OPERATIONS

	if err := c.user.Update(ctx, claims, userID, uu, now); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return nil
}

// Delete removes a user from the database.
func (c Core) Delete(ctx context.Context, claims auth.Claims, userID string) error {

	// PERFORM PRE BUSINESS OPERATIONS

	if err := c.user.Delete(ctx, claims, userID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return nil
}

// FindAll retrieves a list of existing users from the database.
func (c Core) FindAll(ctx context.Context) ([]dto.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	users, err := c.user.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return users, nil
}

// FindByID gets the specified user from the database.
func (c Core) FindByID(ctx context.Context, claims auth.Claims, userID string) (dto.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	usr, err := c.user.FindByID(ctx, claims, userID)
	if err != nil {
		return dto.User{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}

// FindByEmail gets the specified user from the database by email.
func (c Core) FindByEmail(ctx context.Context, claims auth.Claims, email string) (dto.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	usr, err := c.user.FindByEmail(ctx, claims, email)
	if err != nil {
		return dto.User{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (c Core) Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	claims, err := c.user.Authenticate(ctx, now, email, password)
	if err != nil {
		return auth.Claims{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return claims, nil
}
