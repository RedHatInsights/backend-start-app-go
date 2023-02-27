##@ Building

SRC_GO := $(shell find . -name \*.go -print)
SRC_SQL := $(shell find . -name \*.sql -print)
SRC_YAML := $(shell find . -name \*.yaml -print)

PACKAGE_BASE = github.com/RHEnVision/provisioning-backend/internal
LDFLAGS = "-X $(PACKAGE_BASE)/version.BuildCommit=$(shell git rev-parse --short HEAD) -X $(PACKAGE_BASE)/version.BuildTime=$(shell date +'%Y-%m-%d_%T')"

.PHONY: run
run: api ## Build and run backend API
	./api

build: api ## Build all binaries

.PHONY: strip
strip: build ## Strip debug information
	strip api

all-deps: $(SRC_GO) $(SRC_SQL) $(SRC_YAML)

api: all-deps ## Build backend API service
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o api ./cmd/api

.PHONY: clean
clean: ## Clean build artifacts
	-rm api