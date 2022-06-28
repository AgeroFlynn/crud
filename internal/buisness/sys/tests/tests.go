// Package tests contains supporting code for running tests.
package tests

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/AgeroFlynn/crud/internal/buisness/repository/dbschema"
	"github.com/AgeroFlynn/crud/internal/buisness/repository/store/user"
	"github.com/AgeroFlynn/crud/internal/buisness/sys/auth"
	"github.com/AgeroFlynn/crud/internal/foundation/docker"
	"github.com/AgeroFlynn/crud/internal/foundation/keystore"
	"github.com/AgeroFlynn/crud/internal/foundation/logger"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
	"io"
	"os"
	"testing"
	"time"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// DBContainer provides configuration for a container to run.
type DBContainer struct {
	Image string
	Tag   string
	Port  string
	Args  []string
}

// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewUnit(t *testing.T, dbc DBContainer) (*zap.SugaredLogger, *pg.DB, func()) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	pool, resource := docker.StartContainer(t, dbc.Image, dbc.Tag, dbc.Port, dbc.Args...)

	pgOpt := &pg.Options{
		User:     "postgres",
		Password: "postgres",
		Addr:     resource.GetHostPort(dbc.Port),
		Database: "postgres",
	}

	t.Log("Waiting for database to be ready ...")

	err := docker.WaitPostgres(t, pgOpt, dbc.Port, pool, resource)
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	db := pg.Connect(pgOpt)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbschema.Migrate(ctx, db); err != nil {
		docker.DumpContainerLogs(t, resource)
		docker.StopContainer(t, resource)
		t.Fatalf("Migrating error: %s", err)
	}

	if err := dbschema.Seed(ctx, db); err != nil {
		docker.DumpContainerLogs(t, resource)
		docker.StopContainer(t, resource)
		t.Fatalf("Seeding error: %s", err)
	}

	log, err := logger.New("TEST")
	if err != nil {
		t.Fatalf("logger error: %s", err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
		docker.StopContainer(t, resource)

		log.Sync()

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old
		//TODO Check logs from ardan!!!
		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}

	return log, db, teardown
}

// Test owns state for running and shutting down tests.
type Test struct {
	DB       *pg.DB
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	Teardown func()

	t *testing.T
}

// NewIntegration creates a database, seeds it, constructs an authenticator.
func NewIntegration(t *testing.T, dbc DBContainer) *Test {
	log, db, teardown := NewUnit(t, dbc)

	// Create RSA keys to enable authentication in our service.
	keyID := "4754d86b-7a6d-4df5-9c65-224741361492"
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// Build an authenticator using this private key and id for the key store.
	auth, err := auth.New(keyID, keystore.NewMap(map[string]*rsa.PrivateKey{keyID: privateKey}))
	if err != nil {
		t.Fatal(err)
	}

	test := Test{
		DB:       db,
		Log:      log,
		Auth:     auth,
		t:        t,
		Teardown: teardown,
	}

	return &test
}

// Token generates an authenticated token for a user.
func (test *Test) Token(email, pass string) string {
	test.t.Log("Generating token for test ...")

	store := user.NewStore(test.Log, test.DB)
	claims, err := store.Authenticate(context.Background(), time.Now(), email, pass)
	if err != nil {
		test.t.Fatal(err)
	}

	token, err := test.Auth.GenerateToken(claims)
	if err != nil {
		test.t.Fatal(err)
	}

	return token
}

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}
