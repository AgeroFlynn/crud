package mid

import (
	"context"
	"errors"
	"fmt"
	"github.com/AgeroFlynn/crud/internal/buisness/sys/auth"
	"github.com/AgeroFlynn/crud/internal/buisness/sys/validate"
	web2 "github.com/AgeroFlynn/crud/internal/foundation/web"
	"net/http"
	"strings"
)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(a *auth.Auth) web2.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web2.Handler) web2.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Expecting: bearer <token>
			authStr := r.Header.Get("authorization")

			// Parse the authorization header.
			parts := strings.Split(authStr, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: bearer <token>")
				return validate.NewRequestError(err, http.StatusUnauthorized)
			}

			// Validate the token is signed by us.
			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return validate.NewRequestError(err, http.StatusUnauthorized)
			}

			// Add claims to the context, so they can be retrieved later.
			ctx = auth.SetClaims(ctx, claims)

			// Call the next handler.
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// Authorize validates that an authenticated user has at least one role from a
// specified list. This method constructs the actual function that is used.
func Authorize(roles ...string) web2.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web2.Handler) web2.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value return failure.
			claims, err := auth.GetClaims(ctx)
			if err != nil {
				return validate.NewRequestError(
					fmt.Errorf("you are not authorized for that action, no claims"),
					http.StatusForbidden,
				)
			}

			if !claims.Authorized(roles...) {
				return validate.NewRequestError(
					fmt.Errorf("you are not authorized for that action, claims[%v] roles[%v]", claims.Roles, roles),
					http.StatusForbidden,
				)
			}

			// Call the next handler.
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
