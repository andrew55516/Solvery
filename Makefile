postgres:
	docker run --name ps_solvery -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pwd123 -d postgres:15.2-alpine

createdb:
	winpty docker exec -it ps_solvery createdb --username=root --owner=root solvery

dropdb:
	winpty docker exec -it ps_solvery dropdb solvery

migrateup:
	migrate -path db/migrations -database "postgresql://root:pwd123@localhost:5432/solvery?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:pwd123@localhost:5432/solvery?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package=mockdb -destination db/mock/store.go Solvery/db/sqlc Store

server:
	go run cmd/main.go

dockerserver:
	docker-compose up

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test mock gocontainer dockerserver server