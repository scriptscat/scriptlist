
check-cago:
ifneq ($(which cago),)
	go install github.com/codfrm/cago/cmd/cago@latest
endif

check-mockgen:
ifneq ($(which mockgen),)
	go install go.uber.org/mock/mockgen@latest
endif

check-golangci-lint:
ifneq ($(which golangci-lint),)
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
endif

swagger: check-cago
	cago gen swag

lint: check-golangci-lint
	golangci-lint run

lint-fix: check-golangci-lint
	golangci-lint run --fix

test: lint
	go test -v ./...

coverage.out cover:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

html-cover: coverage.out
	go tool cover -html=coverage.out
	go tool cover -func=coverage.out

generate: check-mockgen swagger
	go generate ./... -x

cache_proxy:
	GOOS=linux GOARCH=amd64 go build -o bin/cache_proxy cmd/cache_proxy/main.go

goconvey:
	goconvey
