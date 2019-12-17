test:
	go test ./pkg/...

fmt:
	gofmt -w .

proto:
	protoc -I protos protos/ast.proto --go_out=plugins=grpc:pkg/ast

download:
	@echo Download go.mod dependencies
	go mod download

install-tools: download
	@echo Installing tools from tools.go
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: fmt proto test download install-tools
