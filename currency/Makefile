.PHONY: protos



protos:
	protoc -I protos/ protos/currency.proto \
--go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative



	# protoc -I protos/ protos/currency.proto --go_out=plugins=grpc:protos/currency

	# protoc -I protos/ protos/currency.proto --go_out=plugins=grpc:protos/currency

run:
	go run ./cmd/currency/main.go
