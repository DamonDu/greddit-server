export BIN = ./bin/api
export CMD_DIR = ./cmd/api

all: build start

gen:
	go generate ./...

build: gen
	go build -o ${BIN} ${CMD_DIR}

start:
	export $$(cat .env | grep -v ^\# | xargs) && ${BIN}