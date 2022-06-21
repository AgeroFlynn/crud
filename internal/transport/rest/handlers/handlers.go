package handlers

import (
	"github.com/AgeroFlynn/crud/foundation/web"
	"github.com/AgeroFlynn/crud/internal/buisness/web/mid"
	"github.com/AgeroFlynn/crud/internal/transport/rest/handlers/v1/testgrp"
	"go.uber.org/zap"
	"os"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
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
	app.Handle(version, "/test", tgh.Test)
}
