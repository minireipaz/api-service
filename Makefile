.PHONY: openapi_http lint fmt test

include .env
export

openapi_http:
	@echo "Generating OpenAPI documentation..."
	@./scripts/openapi-http.sh

lint:
	@./scripts/lint.sh

fmt:
	goimports -l -w -d -v ./

test:
	@./scripts/test.sh .env
	@./scripts/test.sh .e2e.env
