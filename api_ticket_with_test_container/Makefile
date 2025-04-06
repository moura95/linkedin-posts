include .env

migrate-up:
	migrate -database ${DB_SOURCE} -path db/migrations up

migrate-down:
	migrate -database ${DB_SOURCE} -path db/migrations down --all

migrate-create:
	@read -p "name of migration: " name; \
	migrate create -dir db/migrations -ext sql -seq $$name


down:
	docker compose down --volumes && docker volume prune -f

up:
	docker compose up -d
	sleep 5
	make migrate-up

sqlc:
	sqlc generate

run:
	go run cmd/main.go

start:
	make up
	sleep 5
	make migrate-up
	go run cmd/main.go

restart:
	make down
	make up
	sleep 10
	make migrate-up
	go run cmd/main.go
swag:
	swag init -g cmd/main.go

test:
	go test -v ./...


.PHONY: migrate-up migrate-down migrate-create down up sqlc start run restart swag test sqlc