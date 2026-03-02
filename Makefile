.PHONY: install
install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: gen
gen: gen-go-proto gen-dart-proto

.PHONY: gen-go-proto
gen-go-proto:
	@for proto in ./api/proto/*.proto; do \
		name=$$(basename $$proto .proto); \
		mkdir -p ./api/pb/$${name}pb; \
		protoc --proto_path=./api/proto \
			--go_out=paths=source_relative:./api/pb/$${name}pb \
			--go-grpc_out=paths=source_relative:./api/pb/$${name}pb \
			$$proto; \
	done

	protoc --proto_path=./llm-runner \
		--go_opt=Mllmrunner.proto=github.com/magomedcoder/gen/api/pb/llmrunner \
		--go-grpc_opt=Mllmrunner.proto=github.com/magomedcoder/gen/api/pb/llmrunner \
		--go_out=module=github.com/magomedcoder/gen:. \
		--go-grpc_out=module=github.com/magomedcoder/gen:. \
		./llm-runner/llmrunner.proto

.PHONY: gen-dart-proto
gen-dart-proto:
	mkdir -p ./client-app/lib/generated/grpc_pb
	protoc --proto_path=./api/proto \
		--dart_out=grpc:./client-app/lib/generated/grpc_pb \
		./api/proto/*.proto

.PHONY: run
run:
	go run ./cmd/gen
