.PHONY: openapi_http lint fmt test

include .env
export

force_upload_datasources_with_fixtures:
	@./data/tinybird/scripts/force_upload_datasources_with_fixtures.sh

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
