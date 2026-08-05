DB_SOURCE=postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5432/simple_bank?sslmode=disable

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up
	
DB_SOURCE=postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5432/simple_bank?sslmode=disable

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

sqlc:
	sqlc generate

.PHONY: migrateup migratedown sqlc