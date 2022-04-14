
SUFFIX=
ifeq ($(GOOS),windows)
	SUFFIX=.exe
endif

swagger:
	swag fmt -g internal/interfaces/api/apis.go
	swag init -g internal/interfaces/api/apis.go --parseDependency --parseDepth 2

linux: generate
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o scriptlist ./cmd/app

build: generate
	go build -o scriptlist$(SUFFIX) ./cmd/app

generate: swagger
	go generate ./... -x

test:
	go test -v ./...

wasm:
	GOOS=js GOARCH=wasm go build -o scriptlist.wasm ./cmd/wasm
