DB_URL = "secret"

postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_PASSWORD=MyPassword1! -d postgres

createDb:
	docker exec -it postgres createdb --username=postgres --owner=postgres simple_bank

dropDb:
	docker exec -it postgres dropdb --username=postgres --owner=postgres simple_bank

migrateUp:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateUp1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migrateDown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateDown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go  github.com/nhatnhanchiha/simpleBank/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: postgres createDb dropDb migrateUp migrateDown sqlc test server mock migrateUp1 migrateDown1 db_docs db_schema proto evans