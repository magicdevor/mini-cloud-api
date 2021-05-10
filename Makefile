createdb:
	docker exec -it postgres_alpine createdb --username=root --owner=root mini_db
dropdb:
	docker exec -it postgres_alpine dropdb mini_db
migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/mini_db?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/mini_db?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: createdb dropdb migrateup migratedown sqlc test server