wire:
	wire gen ./internal/app/

lint:
	golangci-lint run --out-format=colored-line-number

# oapi-codegen -generate types,server -package api openapi/openapi.yaml > internal/oapi/api.gen.go
gen:
	oapi-codegen -generate gin-server -package api openapi/openapi.yaml > internal/oapi/api/api.gen.go