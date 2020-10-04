BIN_PATH=bin
BIN_NAME=corona

test:
	go test -v ./... 

build:
	@go build -o $(BIN_PATH)/$(BIN_NAME) corona.go

# TODO: add args support
run: build
	@./$(BIN_PATH)/$(BIN_NAME)

