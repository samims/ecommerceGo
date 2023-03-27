.PHONY: protos




BINARY_NAME=currency-out

all: clean build run

clean:
	rm -f $(BINARY_NAME)

build:
	go build -o $(BINARY_NAME) ./cmd/currency


run: build
	./$(BINARY_NAME) &
	@echo "binary is starting.."
	@echo "Press CTRL+C to stop."
	@trap "make clean" EXIT; \
		read


protos:
	mkdir -p protos/currency
	protoc -I protos/ protos/currency.proto \
	--go_out=protos/currency --go_opt=paths=source_relative \
	--go-grpc_out=protos/currency --go-grpc_opt=paths=source_relative


#kill:
#	@read -p "Enter the port to kill the process running on: " port; \
#	PID=$$(lsof -t -i :$$port); \
#	if [ -z "$$PID" ]; then \
#		echo "No process found running on port $$port"; \
#	else \
#		echo "Killing process $$PID running on port $$port..."; \
#		kill -9 $$PID; \
#	fi

kill:
	@read -p "Enter the port to kill the process running on: " port; \
	PID=$$(lsof -t -i :$$port); \
	if [ -z "$$PID" ]; then \
		echo "No process found running on port $$port"; \
	else \
		echo "Killing process $$PID running on port $$port..."; \
		kill -9 $$PID; \
	fi