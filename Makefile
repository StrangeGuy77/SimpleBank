DB_SOURCE=postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5432/simple_bank?sslmode=disable

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up
	
migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

sqlc:
	sqlc generate

test:
	go test -count=1 -v -cover ./...

.PHONY: migrateup migratedown sqlc test