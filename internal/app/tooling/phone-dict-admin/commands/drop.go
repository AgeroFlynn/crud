package commands

import (
	"fmt"
	"github.com/AgeroFlynn/crud/internal/buisness/repository/dbschema"
	"github.com/AgeroFlynn/crud/internal/foundation/database"
	"github.com/go-pg/pg/v10"
)

// Drop deletes schema and data from database.
func Drop(opt *pg.Options) error {
	db, err := database.NewPostgresConnection(opt)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	if err := dbschema.DeleteAll(db); err != nil {
		return fmt.Errorf("drop database: %w", err)
	}

	fmt.Println("drop data complete")
	return nil
}
