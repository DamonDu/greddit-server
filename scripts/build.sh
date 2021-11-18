set -eux
export GOOS=linux
export CGO_ENABLED=1
export CC=$(pwd)/scripts/zcc.sh
export CXX=$(pwd)/scripts/zxx.sh

GOARCH=amd64 \
ZTARGET=x86_64-linux-musl \
go build -ldflags="-linkmode external" -o bin/server ./cmd/server/main.go