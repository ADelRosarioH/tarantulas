docker run --rm -v "$PWD":/usr/src/tarantulas -w /usr/src/tarantulas/cmd/tarantulas -e GOOS=linux -e GOARCH=amd64 golang:1.16 go build -v