GITROOT ?=  $(shell git rev-parse --show-toplevel)
BUILD_DIR ?= $(GITROOT)/build

$(BIN_DIR):
	mkdir -p $@

# build binary
## clean: clean all built binary file
clean:
	@rm -rf ${BUILD_DIR}
	@echo "clean all binary file"

build-linux: clean
	GOOS=linux go build -o ${BUILD_DIR}/api -v ./cmd/api/main.go

up-local: build-linux
	@echo "deploy api"
	@docker-compose -f docker-compose.yaml up --build -d

down-local:
	@echo "stop api"
	@docker-compose -f docker-compose.yaml down