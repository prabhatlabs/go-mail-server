include ./prod.env
export

MIGRATIONS_DIR=./db/migrations

# database
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)         # name="something" <- migration name

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1  # this 1 means, rollback 1 migration

db-generate:
	sqlc generate -f ./sqlc.yaml

run:
	go run ./cmd/main.go

build:
	go build -o gomailserver ./cmd/main.go

start:
	./gomailserver
