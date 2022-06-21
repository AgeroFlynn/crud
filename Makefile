
# ==============================================================================
# Local build

run:
	go  run internal/app/services/phone-dict-api/main.go | go run internal/app/tooling/logfmt/main.go

build:
	go build -ldflags "-X main.build=local"

admin:
	go run app/tooling/sales-admin/main.go

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
.PHONY: all test clean phone-dict