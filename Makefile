postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	go run ./db/migration/migrate.go
# migratedown:
# 	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/s14t284/simplebank/db/ent Store

ent:
	go generate ./ent
startdb:
	docker container start postgres12

.PHONY: createdb, dropdb, postgres, migrateup, test, server, mock, ent, startdb
