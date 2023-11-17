db:
	docker compose up -d 
restartdb:
	docker compose down
	docker compose up -d

migrateup:
	migrate -path db/migration -database "postgresql://root:root@127.0.0.1:5434/simple_bank?sslmode=disable" -verbose up           

migratedown:
	migrate -path db/migration -database "postgresql://root:root@127.0.0.1:5434/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:root@127.0.0.1:5434/simple_bank?sslmode=disable" -verbose up 1           

migratedown1:
	migrate -path db/migration -database "postgresql://root:root@127.0.0.1:5434/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/biskitsx/go-backend-master-class/db/sqlc Store

.PHONY: db migrateup migratedown sqlc server