postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=MyPassword1! -d postgres

createDb:
	docker exec -it postgres createdb --username=postgres --owner=postgres simple_bank

dropDb:
	docker exec -it postgres dropdb --username=postgres --owner=postgres simple_bank

migrateUp:
	migrate -path db/migration -database "postgresql://postgres:MyPassword1!@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateUp1:
	migrate -path db/migration -database "postgresql://postgres:MyPassword1!@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrateDown:
	migrate -path db/migration -database "postgresql://postgres:MyPassword1!@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateDown1:
	migrate -path db/migration -database "postgresql://postgres:MyPassword1!@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go  github.com/nhatnhanchiha/simpleBank/db/sqlc Store

.PHONY: postgres createDb dropDb migrateUp migrateDown sqlc test server mock migrateUp1 migrateDown1