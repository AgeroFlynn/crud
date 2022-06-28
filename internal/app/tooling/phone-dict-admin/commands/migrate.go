package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/AgeroFlynn/crud/internal/buisness/repository/dbschema"
	"github.com/AgeroFlynn/crud/internal/foundation/database"
	"github.com/go-pg/pg/v10"
	"time"
)

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

// Migrate creates the schema in the database.
func Migrate(opt *pg.Options) error {
	db, err := database.NewPostgresConnection(opt)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbschema.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migrations complete")
	return nil
}
