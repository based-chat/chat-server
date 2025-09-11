.PHONY: all install-deps generate generate-chat-api install-golangci-lint lint lint-feature clean test
all: clean generate lint check-coverage

test: 
	go test -v -race ./...

clean:
	rm -rf bin coverage.out
	rm -f pkg/chat/v1/*.pb.go pkg/chat/v1/*_grpc.pb.go
	@if [ -d pkg/chat/v1 ] && [ ! "$(ls -A pkg/chat/v1)" ]; then rmdir pkg/chat/v1; fi


LOCAL_BIN?=$(CURDIR)/bin
PROTOC ?= protoc
install-deps:
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.9
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
generate: install-deps
	$(MAKE) generate-chat-api

generate-chat-api: install-deps
	mkdir -p pkg/chat/v1
	@if ! command -v $(PROTOC) >/dev/null 2>&1 ; then echo "$(PROTOC) not found in PATH"; exit 1; fi
	$(PROTOC) \
	--proto_path api/chat/v1 \
	--go_out=pkg/chat/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go-grpc_out=pkg/chat/v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	api/chat/v1/chat.proto

install-golangci-lint:
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0

lint: install-golangci-lint
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.yaml

lint-feature: install-golangci-lint
	$(LOCAL_BIN)/golangci-lint run --config .golangci.yaml --new-from-rev dev

.PHONY: install-go-test-coverage
install-go-test-coverage:
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./coverage.out  -covermode=atomic -coverpkg=./...
	$(LOCAL_BIN)/go-test-coverage --config=./.testcoverage.yml
