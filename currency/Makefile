.PHONY: protos


protos:
	mkdir -p protos/currency
	protoc -I protos/ protos/currency.proto \
	--go_out=protos/currency --go_opt=paths=source_relative \
	--go-grpc_out=protos/currency --go-grpc_opt=paths=source_relative

run:
	go run ./cmd/currency/main.go
