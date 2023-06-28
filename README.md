# grpcmtls
Golang mTLS grpc server and client

1. Generate certificates with [genkeys.sh](./genkeys.sh).
2. Create [grpc API](./api/example.proto) and go generate with `go generate ./...`
3. [Implement](./api/echosvc.go) the API
4. Execute [example](./example_test.go).
