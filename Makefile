
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
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
endif

check-goconvey:
ifneq ($(which goconvey),)
	go install github.com/smartystreets/goconvey@latest
endif


lint: check-golangci-lint
	golangci-lint run

lint-fix: check-golangci-lint
	golangci-lint run --fix

test: lint-fix
	go test -v ./...

coverage.out cover:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

html-cover: coverage.out
	go tool cover -html=coverage.out
	go tool cover -func=coverage.out

generate: check-mockgen check-cago
	go generate ./...
	cago gen

goconvey: check-goconvey
	goconvey

GOOS=linux
GOARCH=amd64
APP_NAME=scriptlist
APP_VERSION=1.0.0

SUFFIX=
ifeq ($(GOOS),windows)
	SUFFIX=.exe
endif

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(APP_NAME)_v$(APP_VERSION)$(SUFFIX) ./cmd/app

cache_proxy:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/cache_proxy$(SUFFIX) cmd/cache_proxy/main.go

