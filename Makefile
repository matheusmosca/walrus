.PHONY: run-server
run-server:
ifeq (, $(wildcard ./build/server))
	@echo "server binary not found, trying to compile it.."
	make compile
endif
	@./build/server

.PHONY: compile
compile:
	@echo "=> installing dependencies..."
	go mod tidy
	@echo "==> Compiling walrus..."
	go build -o build/server cmd/server/v1/main.go

.PHONY: test
test:
	@echo "==> Running tests..."
	go test ./... --race -count=1 -v

.PHONY: lint
lint:
ifeq (, $(shell which $$(go env GOPATH)/bin/golangci-lint))
	@echo "==> golangci-lint not installed, trying to install it..."
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.42.1
endif
	$$(go env GOPATH)/bin/golangci-lint run -c ./.golangci.yml ./...

.PHONY: compile-proto
compile-proto:
	@echo "==> Checking buf dependencies..."
ifeq (, $(shell command -v buf 2> /dev/null))
	@echo "==> Setup: Buf not installed, please follow the instructions on https://docs.buf.build/installation"
endif
	@echo "==> compiling proto..."
	@echo "===> generating grpc code..."
	@go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@buf generate
