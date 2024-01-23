SERVICE=go-transaction-service
BUILD_TIME=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
COMMIT=`git rev-parse HEAD`
LINT_TOOL=$(shell go env GOPATH)/bin/golangci-lint
LINT_VERSION=v1.54.2
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

ifeq "$(shell uname -p)" "arm"
	GO_BUILD_ARCH=arm64
else
	GO_BUILD_ARCH=amd64
endif
ifeq "$(shell uname -s)" "Darwin"
	BUILD_HOST=darwin
else
	BUILD_HOST=linux
endif

SWAGGER_DIR:=$(ROOT_DIR)/swagger
SWAGGER_BIN_DIR:=/usr/local/bin
SWAGGER_TOOL_VERSION=v0.30.5
SWAGGER_ASSET="swagger_$(BUILD_HOST)_$(GO_BUILD_ARCH)"
SWAGGER_ASSET_URL="https://github.com/go-swagger/go-swagger/releases/download/$(SWAGGER_TOOL_VERSION)/$(SWAGGER_ASSET)"

setup-dev: swagger-tool swagger-gin
	@go get -u golang.org/x/tools/cmd/goimports
	@echo "==> Installing linter..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(LINT_VERSION)
	@echo "golangci-lint version:" && golangci-lint version

setup: setup-dev

swagger-gin:
	@echo "==> Installing swagger tool..."	
	@go get -u github.com/swaggo/swag/cmd/swag
	@go get -u github.com/swaggo/gin-swagger
	@go get -u github.com/swaggo/files

swagger-gin-gen:
	@echo "==> Generating swagger..."	
	@swag init -g cmd/app/http/main.go

lint: 
	@echo "==> Linting..."
	$(LINT_TOOL) --version
	$(LINT_TOOL) run --config=.golangci.yaml ./...

test:
	@echo "==> Running unit tests..."
	@go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/)

clean:
	@rm -rf ./bin

deps:
	@echo "==> Tidying module"
	@go mod tidy
	@go mod download

prep:
	@mkdir -p bin/

build-lambda: prep deps swagger-gin-gen
	@echo "==> Building linux arm64 static AWS binary for linux using go $(shell go version)..."
	env GIN_MODE=release CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(SERVICE) -a $(LDFLAGS) ./cmd/app/lambda/main.go	

package-lambda: clean build-lambda
	@echo "==> Packaging lambda for distribution..."
	@cd bin && rm -f $(SERVICE).zip && zip -r $(SERVICE).zip . && cd $(ROOT_DIR)

build-mac: build
build: deps swagger-gin-gen
	@echo "==> Building $(GO_BUILD_ARCH) static binary for $(BUILD_HOST) using go $(shell go version)..."
	env CGO_ENABLED=0 GOOS=$(BUILD_HOST) GOARCH=$(GO_BUILD_ARCH) go build -o bin/$(SERVICE) -a $(LDFLAGS) ./cmd/app/http/main.go	
	chmod +x bin/$(SERVICE)

run-local: build
	@echo "==> Running $(SERVICE) locally..."
	@./bin/$(SERVICE)

deploy:
	@echo "==> Deploying $(SERVICE) to AWS..."
	@sam build && sam deploy --no-confirm-changeset
