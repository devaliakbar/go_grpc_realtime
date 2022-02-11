# Go gRPC realtime demo

A gRPC realtime server demo in Go

## Getting Started

Step 1: Install protoc, see the instructions on
[the Protocol Buffers website](https://developers.google.com/protocol-buffers/).

Step 2: Get the Go protoc-gen by running

```sh
$ go get -u google.golang.org/protobuf/cmd/protoc-gen-go
$ go install google.golang.org/protobuf/cmd/protoc-gen-go

$ go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

Step 3: Add `go/bin/` to your PATH

Step 4: Run this command to generate protoc Go files

```sh
$ protoc --go_out=. --go-grpc_out=. proto/*.proto
```

Step 5: Run the app by `go run lib/main.go` 