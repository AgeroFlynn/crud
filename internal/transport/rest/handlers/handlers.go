package handlers

import (
	userCore "github.com/AgeroFlynn/crud/internal/buisness/core/user"
	"github.com/AgeroFlynn/crud/internal/buisness/sys/auth"
	"github.com/AgeroFlynn/crud/internal/buisness/web/mid"
	"github.com/AgeroFlynn/crud/internal/foundation/web"
	"github.com/AgeroFlynn/crud/internal/transport/rest/handlers/v1/testgrp"
	"github.com/AgeroFlynn/crud/internal/transport/rest/handlers/v1/usergrp"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
	"net/http"
	"os"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	DB       *pg.DB
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {

	// Construct the web.App which holds all routes as well as common Middleware.
	mux := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Panics(),
	)

	// Load the v1 routes.
	v1(mux, cfg)

	return mux
}

// v1 binds all the version 1 routes.
func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
	app.Handle(http.MethodGet, version, "/test", tgh.Test)

	// Register user management and authentication endpoints.
	ugh := usergrp.Handlers{
		User: userCore.NewCore(cfg.Log, cfg.DB),
		Auth: cfg.Auth,
	}

	app.Handle(http.MethodGet, version, "/users/token", ugh.Token)
	app.Handle(http.MethodGet, version, "/users", ugh.FindAll, mid.Authenticate(cfg.Auth), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodGet, version, "/users/{id}", ugh.FindByID, mid.Authenticate(cfg.Auth))
	app.Handle(http.MethodPost, version, "/users", ugh.Create, mid.Authenticate(cfg.Auth), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodPut, version, "/users/{id}", ugh.Update, mid.Authenticate(cfg.Auth), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodDelete, version, "/users/{id}", ugh.Delete, mid.Authenticate(cfg.Auth), mid.Authorize(auth.RoleAdmin))
}
