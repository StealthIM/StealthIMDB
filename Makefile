PROTOCCMD = protoc
PROTOGEN_PATH = $(shell which protoc-gen-go) 
PROTOGENGRPC_PATH = $(shell which protoc-gen-go-grpc) 

GO_FILES := $(shell find $(SRC_DIR) -name '*.go')

GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean

LDFLAGS := -s -w

ifeq ($(OS), Windows_NT)
	DEFAULT_BUILD_FILENAME := StealthIMDB.exe
else
	DEFAULT_BUILD_FILENAME := StealthIMDB
endif

.PHONY: run
run: build
	./bin/$(DEFAULT_BUILD_FILENAME)

StealthIM.DBGateway/db_gateway_grpc.pb.go StealthIM.DBGateway/db_gateway.pb.go: proto/db_gateway.proto
	$(PROTOCCMD) --plugin=protoc-gen-go=$(PROTOGEN_PATH) --plugin=protoc-gen-go-grpc=$(PROTOGENGRPC_PATH) --go-grpc_out=. --go_out=. proto/db_gateway.proto

.PHONY: proto
proto: ./StealthIM.DBGateway/db_gateway_grpc.pb.go ./StealthIM.DBGateway/db_gateway.pb.go

.PHONY: build
build: ./bin/$(DEFAULT_BUILD_FILENAME)

./bin/StealthIMDB.exe: $(GO_FILES) proto
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/StealthIMDB.exe

./bin/StealthIMDB: $(GO_FILES) proto
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/StealthIMDB

.PHONY: build_win build_linux
build_win: ./bin/StealthIMDB.exe
build_linux: ./bin/StealthIMDB

.PHONY: docker_run
docker_run:
	docker-compose up

./bin/StealthIMDB.docker.zst: $(GO_FILES) proto
	docker-compose build
	docker save stealthimdb-app > ./bin/StealthIMDB.docker
	zstd ./bin/StealthIMDB.docker -19
	@rm ./bin/StealthIMDB.docker

.PHONY: build_docker
build_docker: ./bin/StealthIMDB.docker.zst

.PHONY: release
release: build_win build_linux build_docker

.PHONY: clean
clean:
	@rm -rf ./StealthIM.DBGateway
	@rm -rf ./bin
	@rm -rf ./__debug*

.PHONY: debug_proto
debug_proto:
	cd test && python -m grpc_tools.protoc -I. --python_out=. --mypy_out=.  --grpclib_python_out=. --proto_path=../proto db_gateway.proto

