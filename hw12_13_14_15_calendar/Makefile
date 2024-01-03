BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
APP_IMG="calendar-app"
SCHEDULER_IMG="scheduler-app"
POSTGRES_IMG="calendar-db"
RABBITMQ_IMG="rb"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN_CALENDAR) --config=./configs/config.yaml
	#$(BIN_SCHEDULER) --config=./configs/scheduler_config.yaml
	#$(BIN_SENDER) --config=./configs/sender_config.yaml

build-img:
	docker rm -f $(APP_IMG) && docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(APP_IMG) \
		-f build/Dockerfile .

build-postgres:
	docker rm -f $(POSTGRES_IMG) && docker build \
		-t $(POSTGRES_IMG) \
		-f build/postgres/Dockerfile .

run-img: build-img
	docker run -p 8080:8080 -p 9090:9090 --name=$(APP_IMG) $(APP_IMG)

run-postgres: build-postgres
	docker run -d -v pgdata:/var/lib/postgresql/data -p 5432:5432 --name=$(POSTGRES_IMG) $(POSTGRES_IMG)

run-rabbitmq:
	docker rm -f $(RABBITMQ_IMG) && docker run -d --name rb -p 15672:15672 -p 5672:5672 rabbitmq:3-management

version: build
	$(BIN_CALENDAR) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.52.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

migrate:
	goose -dir="./migrations" postgres "$(CONFIG_DATABASE_DSN)" up

generate:
	buf generate
