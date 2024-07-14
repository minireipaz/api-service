include .env
export

.PHONY: openapi_http
openapi_http:
	@echo "Generating OpenAPI documentation..."
	@./scripts/openapi-http.sh api main

.PHONY: lint
lint:
	@./scripts/lint.sh

.PHONY: fmt
fmt:
	goimports -l -w -d -v ./

test:
	@./scripts/test.sh .env
	@./scripts/test.sh .e2e.env
