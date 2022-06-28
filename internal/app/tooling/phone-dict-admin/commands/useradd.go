package commands

import (
	"context"
	"fmt"
	"github.com/AgeroFlynn/crud/internal/buisness/core/dto"
	"github.com/AgeroFlynn/crud/internal/buisness/repository/store/user"
	"github.com/AgeroFlynn/crud/internal/buisness/sys/auth"
	"github.com/AgeroFlynn/crud/internal/foundation/database"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
	"time"
)

// UserAdd adds new users into the database.
func UserAdd(log *zap.SugaredLogger, opt *pg.Options, name, email, password string) error {
	if name == "" || email == "" || password == "" {
		fmt.Println("help: useradd <name> <email> <password>")
		return ErrHelp
	}

	db, err := database.NewPostgresConnection(opt)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	store := user.NewStore(log, db)

	nu := dto.NewUser{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{auth.RoleAdmin, auth.RoleUser},
	}

	usr, err := store.Create(ctx, nu, time.Now())
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
