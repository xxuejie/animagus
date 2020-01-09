test:
	go test ./...

fmt:
	gofmt -w .

proto:
	protoc -I protos protos/ast.proto --go_out=plugins=grpc,paths=source_relative:./pkg/ast
	protoc -I protos protos/generic.proto --go_out=plugins=grpc,paths=source_relative:./pkg/generic

download:
	@echo Download go.mod dependencies
	go mod download

install-tools: download
	@echo Installing tools from tools.go
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

ci-fmt: fmt
	git diff --exit-code

.PHONY: fmt proto test download install-tools ci-fmt
