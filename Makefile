.PHONY: build deploy clean gen

build: clean
	export GO111MODULE=on
	export CGO_ENABLED=1
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/server ./cmd/server/main.go

lambda: clean
	scripts/lambda.sh

deploy: clean lambda
	sls deploy --verbose

clean:
	rm -rf ./bin

gen:
	go generate ./...