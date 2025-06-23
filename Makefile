run:
	go run cmd/server/main.go
test:
	go test -v -count=1 ./...

test-api:
	go test -v -count=1 ./internal/api/router/