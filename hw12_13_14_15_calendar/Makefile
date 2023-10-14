BIN := "./bin/calendar"
APP_IMG="calendar-app"
POSTGRES_IMG="calendar-db"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) --config=./configs/config.yaml

build-img:	build-postgres
	docker rm -f $(APP_IMG) && docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(APP_IMG) \
		-f build/Dockerfile .

build-postgres:
	docker rm -f $(POSTGRES_IMG) && docker build \
		-t $(POSTGRES_IMG) \
		-f build/postgres/Dockerfile .

run-img: build-img build-postgres run-postgres
	docker run -p 8080:8080 --name=$(APP_IMG) $(APP_IMG)

run-postgres:
	docker run -d -v pgdata:/var/lib/postgresql/data -p 5432:5432 --name=$(POSTGRES_IMG) $(POSTGRES_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.52.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

migrate:
	goose -dir="./migrations" postgres "$(CONFIG_DATABASE_DSN)" up

