// Package database provides support for access the database.
package database

import (
	"context"
	"errors"
	"time"
)

import "github.com/go-pg/pg/v10"

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("pg: no rows in result set")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrForbidden             = errors.New("attempted action is not allowed")
)

func NewPostgresConnection(options *pg.Options) (*pg.DB, error) {
	db := pg.Connect(options)

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}
	return db, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *pg.DB) error {

	// First check we can ping the database.
	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.Ping(ctx)
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	// Make sure we didn't timeout or be cancelled.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity. Running this query forces a
	// round trip through the database.
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	return err
}
