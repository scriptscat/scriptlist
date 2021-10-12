
linux:
	GOOS=linux go build -o scriptweb ./cmd/app/main.go

test:
	go test -v ./...
