# Run locally (with Swagger docs)
run:
	go run cmd/server/main.go

# Run in Docker (see docker-build and docker-run)

docker-build:
	docker build -t license-api .

docker-run:
	docker run --rm --env-file .env -p 8080:8080 --name license-api license-api

docker-stop:
	docker stop license-api || true

test:
	go test -v -count=1 ./...

test-api:
	go test -v -count=1 ./internal/api/router/

test-storage:
	go test -v -count=1 ./internal/storage/


swagger:
	swag init -q -g cmd/server/main.go
