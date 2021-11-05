PROJECT_NAME=restapi
VERSION = 3.0
CI_PIPELINE_IID ?= 0
BUILD_VERSION = $(VERSION).$(CI_PIPELINE_IID)
BRANCH = $(CI_COMMIT_REF_NAME)
CURRENT_DIR?=$(shell pwd)

version:
	@echo $(BUILD_VERSION)

build: build-docs
	docker build --build-arg VERSION=$(BUILD_VERSION) --build-arg BRANCH=$(BRANCH) --build-arg http_proxy=$(HTTP_PROXY) . -t $(PROJECT_NAME):$(BUILD_VERSION)

build-docs:
	@docker run --name $(PROJECT_NAME)-docs --rm -i -v "$(CURRENT_DIR):/restapi" -e "GOPATH=/restapi" -w /restapi openapi-godocgen:latest sh -c "openapi-godocgen ./src/restapi/handlers/ > ./docs/internal/index.yml"

up: migrate-up
	docker-compose up restapi

test:
	@docker run --name $(PROJECT_NAME)-test --rm -i -v "$(CURRENT_DIR):/restapi" -e "GOPATH=/restapi" -w /restapi golang:1.13 bash -c "go test ./src/... -coverprofile=coverage.out > test-report.txt && go tool cover -func=coverage.out >> test-report.txt || cat test-report.txt"
	make build-docs

lint:
	@docker run --name $(PROJECT_NAME)-lint --rm -i -v "$(CURRENT_DIR):/restapi" -e "GOPATH=/restapi" -w /restapi golangci/golangci-lint:v1.23-alpine golangci-lint run ./src/restapi/... -E gofmt --skip-dirs=./src/restapi/vendor --deadline=5m

coverage:
	@docker run --name $(PROJECT_NAME)-coverage --rm -i -v "$(CURRENT_DIR):/restapi" -e "GOPATH=/restapi" -w /restapi test-coverage:latest bash -c "test-coverage -report=test-report.txt -limit=7"

push:
	@docker tag $(PROJECT_NAME):$(BUILD_VERSION) $(REGISTRY)/$(PROJECT_NAME):$(BUILD_VERSION)
	@docker tag $(PROJECT_NAME):$(BUILD_VERSION) $(REGISTRY)/$(PROJECT_NAME):$(VERSION)
	@docker tag $(PROJECT_NAME):$(BUILD_VERSION) $(REGISTRY)/$(PROJECT_NAME):latest
	@docker push $(REGISTRY)/$(PROJECT_NAME)

migrate-up:
	docker-compose up migrate-up

migrate-down:
	docker-compose up migrate-down

vendor:
	export GOPATH=`pwd` && cd ./src/restapi && go mod vendor

tidy:
	export GOPATH=`pwd` && cd ./src/restapi && go mod tidy

protoc:
	@docker run --name $(PROJECT_NAME)-protoc --rm -i -v "$(CURRENT_DIR):/proto" -w /proto onrik/protoc-gen-micro:1.0.0 sh -c "\
	protoc --go_out=./src/restapi/proto/eds --micro_out=./src/restapi/proto/eds -I=proto/eds ./proto/eds/*.proto && \
	protoc --go_out=./src/restapi/proto/ccc --micro_out=./src/restapi/proto/ccc -I=proto/ccc ./proto/ccc/*.proto && \
	protoc --go_out=./src/restapi/proto/ccc --micro_out=./src/restapi/proto/ccc -I=proto/ccc ./proto/ccc/*.proto && \
	protoc --go_out=./src/restapi/proto/dds --micro_out=./src/restapi/proto/dds -I=proto/dds ./proto/dds/*.proto && \
	protoc --go_out=./src/restapi/proto/tip --micro_out=./src/restapi/proto/tip -I=proto/tip ./proto/tip/*.proto && \
	protoc --go_out=./src/restapi/proto/healthboard --micro_out=./src/restapi/proto/healthboard -I=proto/healthboard ./proto/healthboard/*.proto && \
	protoc --go_out=./src/restapi/proto/news --micro_out=./src/restapi/proto/news -I=proto/news ./proto/news/*.proto && \
	protoc --go_out=./src/restapi/proto/github_crawler --micro_out=./src/restapi/proto/github_crawler -I=proto/github_crawler ./proto/github_crawler/*.proto && \
	protoc --go_out=./src/restapi/proto/cve --micro_out=./src/restapi/proto/cve -I=proto/cve ./proto/cve/*.proto && \
	protoc --go_out=./src/restapi/proto/shodan --micro_out=./src/restapi/proto/shodan -I=proto/shodan ./proto/shodan/*.proto && \
	protoc --go_out=./src/restapi/proto/whitelist --micro_out=./src/restapi/proto/whitelist -I=proto/whitelist ./proto/whitelist/*.proto && \
	protoc --go_out=./src/restapi/proto/report_parser --micro_out=./src/restapi/proto/report_parser -I=proto/report_parser ./proto/report_parser/*.proto && \
	protoc --go_out=./src/restapi/proto/retro --micro_out=./src/restapi/proto/retro -I=proto/retro ./proto/retro/*.proto && \
	protoc --go_out=./src/restapi/proto/mitrerulesrepo --micro_out=./src/restapi/proto/mitrerulesrepo -I=proto/mitrerulesrepo ./proto/mitrerulesrepo/*.proto && \
	protoc --go_out=./src/restapi/proto/kyc --micro_out=./src/restapi/proto/kyc -I=proto/kyc ./proto/kyc/*.proto && \
	protoc --go_out=./src/restapi/proto/mitrerulesrepo --micro_out=./src/restapi/proto/mitrerulesrepo -I=proto/mitrerulesrepo ./proto/mitrerulesrepo/*.proto && \
	protoc --go_out=./src/restapi/proto/vulners --micro_out=./src/restapi/proto/vulners -I=proto/vulners ./proto/vulners/*.proto && \
	protoc --go_out=./src/restapi/proto/ioctracker --micro_out=./src/restapi/proto/ioctracker -I=proto/ioctracker ./proto/ioctracker/*.proto && \
	protoc --go_out=./src/restapi/proto/geoip --micro_out=./src/restapi/proto/geoip -I=proto/geoip ./proto/geoip/*.proto && \
	protoc --go_out=./src/restapi/proto/storage --micro_out=./src/restapi/proto/storage -I=proto/storage ./proto/storage/*.proto && \
	protoc --go_out=./src/restapi/proto/pastebin --micro_out=./src/restapi/proto/pastebin -I=proto/pastebin ./proto/pastebin/*.proto"
