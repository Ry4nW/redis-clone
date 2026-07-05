run:
	go run ./cmd/server

test:
	go test ./...

race:
	go test -race ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
