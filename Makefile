BIN_PATH=bin
BIN_NAME=corona

ifeq (run, $(firstword $(MAKECMDGOALS)))  
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))  
  $(eval $(RUN_ARGS):;@:)
endif


test:
	go test -v ./... 

build:
	@go build -o $(BIN_PATH)/$(BIN_NAME) corona.go

.PHONY: run
run: build
	./$(BIN_PATH)/$(BIN_NAME) $(RUN_ARGS)
