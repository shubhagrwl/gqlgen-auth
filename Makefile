lint:
	golangci-lint run

migration:
	@read -p "migration file name:" module; \
	cd /internal/app/db/migrations && goose create $$module sql

generate:
	go run github.com/99designs/gqlgen generate -c internal/app/service