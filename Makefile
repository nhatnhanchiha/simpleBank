postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=MyPassword1! -d postgres

createDb:
	docker exec -it postgres createdb --username=postgres --owner=postgres simple_bank

dropDb:
	docker exec -it postgres dropdb --username=postgres --owner=postgres simple_bank

migrateUp:
	migrate -path db/migration -database "postgresql://postgres:MyPassword1!@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateDown:
	migrate -path db/migration -database "postgresql://postgres:MyPassword1!@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createDb dropDb migrateUp migrateDown sqlc test