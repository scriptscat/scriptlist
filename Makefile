
SUFFIX=
ifeq ($(GOOS),windows)
	SUFFIX=.exe
endif

linux: generate
	CGO_ENABLE=0 GOOS=linux go build -o scriptlist ./cmd/app

build: generate
	go build -o scriptlist$(SUFFIX) ./cmd/app

generate:
	go generate ./... -x

test:
	go test -v ./...

wasm:
	GOOS=js GOARCH=wasm go build -o scriptlist.wasm ./cmd/wasm
