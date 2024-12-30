DB_URL=postgresql://postgres:postgres@localhost:5433/eshop?sslmode=disable
create-migration:
	@echo "Creating migration file..."
	migrate create -ext sql -dir migrations -seq $(name)
migrate-up:
	@echo "Running migrations..."
	migrate -path migrations -database "$(DB_URL)" up
migrate-up-1:
	@echo "Running migrations up 1..."
	migrate -path migrations -database "$(DB_URL)" up 1
migrate-down:
	@echo "Running migrations down..."
	migrate -path migrations -database "$(DB_URL)" down
migrate-down-1:
	@echo "Running migrations down 1..."
	migrate -path migrations -database "$(DB_URL)" down 1
migrate-drop:
	@echo "Running migrations..."
	migrate -path migrations -database "$(DB_URL)" drop
serve:
	@echo "Starting server..."
	go run main.go
sqlc:
	@echo "Generating sqlc..."
	sqlc generate
build:
	@echo "Building application..."
	go build ./cmd/web
run:
	@echo "Running application..."
	go run ./cmd/web --profile api
run-docker:
	@echo "Running docker..."
	docker-compose up -d
gen-swagger:
	@echo "Generating swagger..."
	swag init -d internal/api -g server.go --parseInternal --parseDependency
.PHONY: create-migration migrate-up serve sqlc build run run-docker gen-swagger migrate-up-1 migrate-down migrate-down-1 migrate-drop