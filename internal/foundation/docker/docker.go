// Package docker provides support for starting and stopping docker containers
// for running tests.
package docker

import (
	"database/sql"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os/exec"
	"testing"
	"time"
)

// StartContainer starts the specified container for running tests.
func StartContainer(t *testing.T, image string, tag string, port string, args ...string) (*dockertest.Pool, *dockertest.Resource) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: image,
		Tag:        tag,
		Env:        args,
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	t.Logf("Image:       %s", image)
	t.Logf("ContainerID: %s", resource.Container.ID)
	t.Logf("Host:        %s", resource.GetHostPort(port))

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	return pool, resource
}

// StopContainer stops and removes the specified container.
func StopContainer(t *testing.T, resource *dockertest.Resource) {
	if err := resource.Close(); err != nil {
		t.Fatalf("Could not close resource: %s", err)
	}
}

// DumpContainerLogs outputs logs from the running docker container.
func DumpContainerLogs(t *testing.T, resource *dockertest.Resource) {
	out, err := exec.Command("docker", "logs", resource.Container.ID).CombinedOutput()
	if err != nil {
		t.Fatalf("could not log container: %v", err)
	}
	t.Logf("Logs for %s\n%s:", resource.Container.ID, out)
}

// WaitPostgres waits necessary time for postgres container start
func WaitPostgres(t *testing.T, opt *pg.Options, port string, pool *dockertest.Pool, resource *dockertest.Resource) error {
	hostAndPort := resource.GetHostPort(port)
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		opt.User,
		opt.Password,
		hostAndPort,
		opt.Database,
	)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err := pool.Retry(func() error {
		db, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}
	return nil
}
