##@ Testing

TEST_TAGS?=test

.PHONY: test
test: ## Run unit tests
	go test -tags=$(TEST_TAGS) ./...

.PHONY: test-database
test-database: ## Run integration tests (require database)
	# "go test pkg1 pkg2" would run tests in parallel causing database locks
	go test --count=1 -v -tags=database ./internal/dao/tests
