
linux:
	CGO_ENABLE=0 GOOS=linux go build -o scriptlist ./cmd/app

test:
	go test -v ./...
