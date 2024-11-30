DB_URL=postgresql://postgres:postgres@localhost:5433/eshop?sslmode=disable
create-migration:
	@echo "Creating migration file..."
	migrate create -ext sql -dir migrations -seq $(name)
migrate-up:
	@echo "Running migrations..."
	migrate -path migrations -database "$(DB_URL)" up
serve:
	@echo "Starting server..."
	go run main.go
sqlc:
	@echo "Generating sqlc..."
	sqlc generate
build:
	@echo "Building application..."
	go build ./cmd/eshop
run-docker:
	@echo "Running docker..."
	docker-compose up -d
.PHONY: create-migration migrate-up serve sqlc build