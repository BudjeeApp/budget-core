-include .env

include build/makefiles/osvars.mk

TEST_APP_ENV := TEST
export TEST_DATABASE_URL := postgres://${LOCAL_DB_USER}:${LOCAL_DB_PASSWORD}@${LOCAL_DB_HOST}:${LOCAL_DB_PORT}/${LOCAL_DB_TEST_DATABASE}?sslmode=disable

run_server:
	go run main.go server

seed_database:
	go run main.go seedDatabase

# Database commands
.PHONY: test_db rollback_test_db drop_test_db
test_db:
	DATABASE_URL=$(TEST_DATABASE_URL) dbmate up

rollback_test_db:
	DATABASE_URL=$(TEST_DATABASE_URL) dbmate rollback

drop_test_db:
	DATABASE_URL=$(TEST_DATABASE_URL) dbmate drop

.PHONY: dev_db rollback_dev_db drop_dev_db
dev_db:
	dbmate up

rollback_dev_db:
	dbmate rollback

drop_dev_db:
	dbmate drop

# Test commands
.PHONY: test coverage_report
test coverage.out:
	APP_ENV=$(TEST_APP_ENV) TEST_DATABASE_URL=$(TEST_DATABASE_URL) go test -v -p 1 ./... -coverpkg=./internal/... -coverprofile ./coverage.out

coverage_report: coverage.out
	go tool cover -html=coverage.out

# Tools
TOOLS := github.com/twitchtv/twirp/protoc-gen-twirp \
		 google.golang.org/protobuf/cmd/protoc-gen-go \
		 github.com/golang/mock/mockgen

TOOLS_DIR := tools_module/tools
TOOLS_BIN := $(TOOLS_DIR)/bin

$(TOOLS_DIR):
	mkdir -v -p $@

$(TOOLS_BIN): tools_vendor
	(cd tools_module; GOBIN="$(PWD)/$(TOOLS_BIN)" go install -mod=vendor $(TOOLS))

tools_vendor: tools_module/tools.go clean_vendor
	(cd tools_module; go mod vendor)

.PHONY: clean_vendor
clean_vendor:
	rm -rf vendor
	rm -rf tools_module/vendor

tools: $(TOOLS_BIN)

clean_tools:
	rm -rf "$(TOOLS_DIR)"

## Protoc
PROTOC_VERSION := 21.12
PROTOC_RELEASES_PATH := https://github.com/protocolbuffers/protobuf/releases/download
PROTOC_ZIP := protoc-$(PROTOC_VERSION)-$(PROTOC_PLATFORM).zip
PROTOC_DOWNLOAD := $(PROTOC_RELEASES_PATH)/v$(PROTOC_VERSION)/$(PROTOC_ZIP)
PROTOC := $(TOOLS_BIN)/protoc

PROTO_DIR := ./rpc/**
PROTOBUF_FILES := $(PROTO_DIR)/service.twirp.go $(PROTO_DIR)/service.pb.go

# protoc
$(PROTOC): $(TOOLS_DIR)/$(PROTOC_ZIP)
	unzip -o -d "$(TOOLS_DIR)" $< && touch $@

$(TOOLS_DIR)/$(PROTOC_ZIP):
	curl --location $(PROTOC_DOWNLOAD) --output $@

setup: tools $(PROTOC)

.PHONY: services clean_services
services:
	 PATH=$(TOOLS_BIN):$(PATH) protoc --go_out=. --twirp_out=. rpc/category/service.proto

clean_services:
	rm -f $(PROTOBUF_FILES)

# Mocks
MOCK_DIR := ./internal/category/mocks

.PHONY: clean_mocks
clean_mocks:
	rm -f $(MOCK_DIR)/*

mocks: \
	$(MOCK_DIR)/mock_repository.go

## Mock Sources
$(MOCK_DIR)/mock_repository.go: internal/category/repository/repository.go

# mockgen commands for generation
$(MOCK_DIR)/mock_repository.go :
	$(TOOLS_BIN)/mockgen -source $< -destination $@ -package=mocks