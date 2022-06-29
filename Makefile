
# For testing a simple query on the system. Don't forget to `make seed` first.
# curl --user "admin@example.com:gophers" http://localhost:3000/v1/users/token
# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users

# ==============================================================================
# Local build

run:
	go  run internal/app/services/phone-dict-api/main.go | go run internal/app/tooling/logfmt/main.go

build:
	go build -ldflags "-X main.build=local"

postgres-init:
	docker run --name local-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres -p 5432:5432 -d postgres:13-alpine

postgres-start:
	docker start local-postgres

# ==============================================================================
# Administration

migrate:
	go run internal/app/tooling/phone-dict-admin/main.go migrate | go run internal/app/tooling/logfmt/main.go

seed: migrate
	go run internal/app/tooling/phone-dict-admin/main.go seed | go run internal/app/tooling/logfmt/main.go

drop:
	go run internal/app/tooling/phone-dict-admin/main.go drop | go run internal/app/tooling/logfmt/main.go

reload: drop seed

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
.PHONY: all test clean phone-dict