.PHONY: install gen gen-proto build run test

install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

gen-proto:
	@for proto in ./api/proto/app/*.proto; do \
		name=$$(basename $$proto .proto); \
		mkdir -p ./api/pb/app/$${name}pb; \
		protoc --proto_path=./api/proto/app \
			--go_out=paths=source_relative:./api/pb/app/$${name}pb \
			--go-grpc_out=paths=source_relative:./api/pb/app/$${name}pb \
			$$proto; \
	done

build:
	mkdir -p ./build
	go build -o ./build/gen-server ./cmd/gen

run:
	go run ./cmd/gen

test:
	go test ./... -race -count=1
