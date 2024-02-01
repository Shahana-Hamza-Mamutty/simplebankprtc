postgres:
	docker run --name postgres16 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=shahana7 -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:shahana7@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:shahana7@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock