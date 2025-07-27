.PHONY: help tidy build clean install-mockery refresh-mocks test cover-report air

ifeq (add-migration,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "add-migration"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

ifeq (add-migration-sqlite,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "add-migration-sqlite"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

DATE=$(shell date +%Y%m%d)
GIT_HASH=$(shell git rev-parse --short HEAD)
UNIX_TS=${DATE}-${GIT_HASH}

tidy:
	go mod tidy

build: clean tidy test
	CGO_ENABLED=0 go build -o bin/runner -ldflags="-w -s -extldflags=-static -X 'github.com/utking/spaces/internal.Version=${UNIX_TS}'" cmd/main.go

clean:
	rm -rf bin/

install-mockery:
ifneq ($(shell test -e /usr/local/bin/mockery && echo -n yes),yes)
	echo "Installing mockery before running the tests"
	go install github.com/vektra/mockery/v3@latest
endif

force-mockery-update:
	go install github.com/vektra/mockery/v3@latest

refresh-mocks: install-mockery
	rm -rf mocks; mockery

test: refresh-mocks
	go test -cover -v ./...

cover-report:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

air:
	air run cmd/main.

install-air:
	go install github.com/air-verse/air@latest

update:
	go get -u all

lint:
	go tool revive ./...

install-lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6

lintci:
	golangci-lint run

install-migrate:
	go install -tags mysql github.com/golang-migrate/migrate/v4/cmd/migrate@latest

_new-migration-sqlite:
	migrate create -ext sql -dir migrations/sqlite3 -seq $(RUN_ARGS)

_new-migration:
	migrate create -ext sql -dir migrations/mysql -seq $(RUN_ARGS)

add-migration: _new-migration

add-migration-sqlite: _new-migration-sqlite

migrate-up:
	go run cmd/main.go migrate

migrate-up-sqlite:
	go run cmd/main.go migrate -s

migrate-down:
	go run cmd/main.go migrate --down

migrate-down-sqlite:
	go run cmd/main.go migrate --down -s

views-lint:
	~/.local/bin/djlint views/ --profile=golang

install-views-lint:
	pip install djlint

lint-js:
	npx oxlint@latest

make linters-all: views-lint lint-js lintci