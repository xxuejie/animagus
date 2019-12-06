test:
	go test ./pkg/...

fmt:
	gofmt -w .

proto:
	protoc -I protos protos/ast.proto --go_out=plugins=grpc:pkg/ast

.PHONY: fmt proto test
