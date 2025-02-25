
.PHONY: serve
serve: static/client.wasm static/wasm_exec.js
	go run ./server/

static/client.wasm: client/main.go
	GOOS=js GOARCH=wasm go build -o static/client.wasm ./client

static/wasm_exec.js:
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" static/

.PHONY: generate
generate:
	buf generate

.PHONY: setup
setup:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest