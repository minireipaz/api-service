#!
docker run --rm \
  -v ${PWD}/definitionapi:/local openapitools/openapi-generator-cli generate \
  -i /local/openapi3.0.yml \
  -g go \
  -o /local/out/go
