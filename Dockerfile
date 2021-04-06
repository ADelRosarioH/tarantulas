FROM golang:1.16-alpine as build
COPY . /usr/src/tarantulas/
WORKDIR /usr/src/tarantulas/cmd/tarantulas
RUN GOOS=linux GOARCH=amd64 go build -v

FROM alpine:latest  
COPY --from=build /usr/src/tarantulas/cmd/tarantulas/tarantulas  /usr/local/bin/tarantulas
