check_install:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

kill:
	@read -p "Enter the port to kill the process running on: " port; \
	PID=$$(lsof -t -i :$$port); \
	if [ -z "$$PID" ]; then \
		echo "No process found running on port $$port"; \
	else \
		echo "Killing process $$PID running on port $$port..."; \
		kill -9 $$PID; \
	fi

BINARY_NAME=productapi-out

all: clean build run

clean:
	rm -f $(BINARY_NAME)

build:
	go build -o $(BINARY_NAME) ./


run: build
	./$(BINARY_NAME) &
	@echo "binary is starting.."
	@echo "Press CTRL+C to stop."
	@trap "make clean" EXIT; \
		read
