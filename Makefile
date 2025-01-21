SHELL := /bin/bash
BASEDIR = $(shell pwd)

# 可在make是带入参数进行替换
# usage: SERVICE_NAME=eagle-service make build 
SERVICE_NAME?=eagle-service

# build with version infos
versionDir = "github.com/go-eagle/eagle/pkg/version"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

ldflags="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"

PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

# proto
API_PROTO_FILES=$(shell find api -name *.proto)
API_PROTO_PB_FILES=$(shell find api -name *.pb.go)

# init environment variables
export PATH        := $(shell go env GOPATH)/bin:$(PATH)
export GOPATH      := $(shell go env GOPATH)
export GO111MODULE := on

# make   make all
.PHONY: all
all: lint test build

.PHONY: build
# make build, Build the binary file
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags ${ldflags} -o bin/$(SERVICE_NAME) cmd/server/main.go cmd/server/wire_gen.go
#	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags ${ldflags} -o bin/$(SERVICE_NAME)-consumer cmd/consumer/main.go cmd/consumer/wire_gen.go

.PHONY: run
# make run, run current project
run: wire
	go run cmd/server/main.go cmd/server/wire_gen.go

.PHONY: wire
# make wire, generate wire_gen.go
wire:
	cd cmd/server && wire

.PHONY: fmt
# make fmt
fmt:
	@gofmt -s -w .

.PHONY: golint
# make golint
golint:
	@if ! which golint &>/dev/null; then \
  		echo "Installing golint"; \
  		go get -u golang.org/x/lint/golint; \
  	fi
	@golint -set_exit_status ${PKG_LIST}

.PHONY: lint
# make lint
lint:
	@if ! which golangci-lint &>/dev/null; then \
  		echo "Installing golangci-lint"; \
  		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0; \
  	fi
	${GOPATH}/bin/golangci-lint run ./...

.PHONY: test
# make test
test: vet
	@go test -race -short ${PKG_LIST}

.PHONY: vet
# make vet
vet:
	go vet ./... | grep -v vendor;true

.PHONY: cover
# make cover
cover:
	@go test -short -coverprofile=coverage.out -covermode=atomic ${PKG_LIST}

.PHONY: view-cover
# make view-cover  preview coverage report
view-cover: cover
	go tool cover -html=coverage.out

.PHONY: image
# make image  生成docker镜像, eg: make image GIT_TAG=v1.0.0
image:
	sh deploy/docker_image.sh ${SERVICE_NAME} $(GIT_TAG)

.PHONY: deploy
# make deploy  deploy app to k8s, eg: make deploy GIT_TAG=v1.0.0
deploy:
	sh deploy/deploy.sh ${SERVICE_NAME} $(GIT_TAG)

.PHONY: clean
# make clean
clean:
	@-rm -vrf bin/${SERVICE_NAME}
	@-rm -vrf cover.out
	@-rm -vrf coverage.txt
	@go mod tidy
	@echo "clean finished"

.PHONY: docs
# gen swagger doc
docs:
	@if ! which swag &>/dev/null; then \
  		echo "downloading swag"; \
  		go get -u github.com/swaggo/swag/cmd/swag; \
  	fi
	@swag init
	@mv docs/docs.go api/http
	@mv docs/swagger.json api/http
	@mv docs/swagger.yaml api/http
	@echo "gen-docs done"
	@echo "see docs by: http://localhost:8080/swagger/index.html"

.PHONY: graph
# make graph 生成交互式的可视化Go程序调用图(会在浏览器自动打开)
graph:
	@export GO111MODULE="on"
	@if ! which go-callvis &>/dev/null; then \
  		echo "downloading go-callvis"; \
  		go get -u github.com/ofabry/go-callvis; \
  	fi
	@echo "generating graph"
	@go-callvis github.com/go-eagle/eagle

.PHONY: mockgen
# make mockgen gen mock file
mockgen:
	@echo "downloading mockgen"
	@go get github.com/golang/mock/mockgen
	cd ./internal &&  for file in `egrep -rnl "type.*?interface" ./dao | grep -v "_test" `; do \
		echo $$file ; \
		cd .. && mockgen -destination="./internal/mock/$$file" -source="./internal/$$file" && cd ./internal ; \
	done

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/google/gnostic@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/favadi/protoc-go-inject-tag@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/gogo/protobuf/protoc-gen-gogo@latest
	go install github.com/gogo/protobuf/protoc-gen-gogofast@latest
	go install github.com/gogo/protobuf/protoc-gen-gogofaster@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: proto
# generate proto struct with validate
proto:
	protoc --proto_path=. \
           --proto_path=./third_party \
           --go_out=. --go_opt=paths=source_relative \
           --validate_out=lang=go,paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: grpc
# generate grpc code with remove omitempty from json tag
grpc:
# note: --gogofaster_out full replace --go_out=. --go_opt=paths=source_relative
	@for v in $(API_PROTO_FILES); do \
  		echo "./$$v"; \
		protoc --proto_path=. \
			   --proto_path=./third_party \
				 --gogofast_out=. --gogofast_opt=paths=source_relative \
         --go-gin_out=. --go-gin_opt=paths=source_relative \
         --go-grpc_out=. --go-grpc_opt=paths=source_relative \
			   "./$$v"; \
    done

.PHONY: tag
# add custom tag to pb struct
tag:
	protoc-go-inject-tag -input=$(API_PROTO_PB_FILES)

.PHONY: openapi
# generate openapi
openapi:
	protoc --proto_path=. \
          --proto_path=./third_party \
          --openapi_out=. \
          $(API_PROTO_FILES)

.PHONY: doc
# generate html or markdown doc
doc:
	protoc --proto_path=. \
          --proto_path=./third_party \
	   	  --doc_out=. \
	   	  --doc_opt=html,index.html \
	   	  $(API_PROTO_FILES)

.PHONY: gorm-gen
# generate gen file for gorm
gorm-gen:
	go run cmd/gen/generate.go

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m  %-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := all
