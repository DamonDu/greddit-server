.PHONY: build deploy clean gen

build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/server ./cmd/server/main.go

lambda:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/lambda ./cmd/lambda/main.go

deploy: clean lambda
	sls deploy --verbose

clean:
	rm -rf ./bin

gen:
	go generate ./...