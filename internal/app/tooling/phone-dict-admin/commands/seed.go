package commands

import (
	"context"
	"fmt"
	"github.com/AgeroFlynn/crud/internal/buisness/repository/dbschema"
	"github.com/AgeroFlynn/crud/internal/foundation/database"
	"github.com/go-pg/pg/v10"
	"time"
)

// Seed loads test data into the database.
func Seed(opt *pg.Options) error {
	db, err := database.NewPostgresConnection(opt)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbschema.Seed(ctx, db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("seed data complete")
	return nil
}
