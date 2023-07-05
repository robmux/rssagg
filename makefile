# Go variables
GO          := go
GOFLAGS     :=
GOTESTFLAGS := -v
GOOSE       := goose
EXECUTABLE  := rssagg

include .env

# Target: build
build:
	$(GO) build $(GOFLAGS) -o $(EXECUTABLE)

# Target: run
run: build
	./$(EXECUTABLE)

# Target: test
test:
	$(GO) test $(GOFLAGS) $(GOTESTFLAGS) ./...

# Target: integration-test
integration-test:
	$(GO) test $(GOFLAGS) $(GOTESTFLAGS) -tags=integration ./...

# Target: u-d to update all the dependencies
u-d:
	$(GO) get -u ./...

# Target: goose-up
goose-up:
	$(GOOSE) -dir ./sql/schema postgres $(DB_URL_G) up

# Target: goose-down
goose-down:
	$(GOOSE) -dir ./sql/schema postgres $(DB_URL_G) down

# Target: goose-create <migration-name>
goose-create:
	$(GOOSE) -dir ./sql/schema create $(migration-name) sql

.PHONY: build run test integration-test update-dependencies goose-up goose-down goose-create
