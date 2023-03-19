##@ Generate

.PHONY: generate-openapi
generate-openapi: ## Generate OpenAPI spec
	go run ./cmd/openapi_spec
