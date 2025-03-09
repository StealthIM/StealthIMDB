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

run: build
	./bin/$(DEFAULT_BUILD_FILENAME)

StealthIM.DBGateway/db_gateway_grpc.pb.go StealthIM.DBGateway/db_gateway.pb.go: proto/db_gateway.proto
	$(PROTOCCMD) --plugin=protoc-gen-go=$(PROTOGEN_PATH) --plugin=protoc-gen-go-grpc=$(PROTOGENGRPC_PATH) --go-grpc_out=. --go_out=. proto/db_gateway.proto

proto: ./StealthIM.DBGateway/db_gateway_grpc.pb.go ./StealthIM.DBGateway/db_gateway.pb.go


build: ./bin/$(DEFAULT_BUILD_FILENAME)

./bin/StealthIMDB.exe: $(GO_FILES) proto
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/StealthIMDB.exe

./bin/StealthIMDB: $(GO_FILES) proto
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/StealthIMDB

build_win: ./bin/StealthIMDB.exe
build_linux: ./bin/StealthIMDB

docker_run:
	docker-compose up

./bin/StealthIMDB.docker.zst: $(GO_FILES) proto
	docker-compose build
	docker save stealthimdb-app > ./bin/StealthIMDB.docker
	zstd ./bin/StealthIMDB.docker -19
	@rm ./bin/StealthIMDB.docker

build_docker: ./bin/StealthIMDB.docker.zst

release: build_win build_linux build_docker

clean:
	@rm -rf ./StealthIM.DBGateway
	@rm -rf ./bin
	@rm -rf ./__debug*
