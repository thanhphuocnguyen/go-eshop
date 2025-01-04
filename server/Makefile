# for server
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
build-migrate:
	@echo "Building migration..."
	go build ./cmd/migrate
build-seed:
	@echo "Building seed..."
	go build ./cmd/seed
build-server:
	@echo "Building application..."
	go build ./cmd/web
serve-server:
	@echo "Running application..."
	go run ./cmd/web --profile api
gen-sqlc:
	@echo "Generating sqlc..."
	sqlc generate
gen-swagger:
	@echo "Generating swagger..."
	swag init -d internal/api -g server.go --parseInternal --parseDependency

.PHONY: create-migration migrate-up migrate-up-1 migrate-down migrate-down-1 migrate-drop 
	build-migrate build-server serve-server gen-sqlc gen-swagger build-seed