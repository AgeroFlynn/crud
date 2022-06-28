package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AgeroFlynn/crud/internal/buisness/repository/store/user"
	"github.com/AgeroFlynn/crud/internal/foundation/database"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
	"os"
	"time"
)

// Users retrieves all users from the database.
func Users(log *zap.SugaredLogger, opt *pg.Options) error {
	db, err := database.NewPostgresConnection(opt)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	store := user.NewStore(log, db)

	users, err := store.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("retrieve users: %w", err)
	}

	return json.NewEncoder(os.Stdout).Encode(users)
}
