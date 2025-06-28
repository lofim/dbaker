export CGO_ENABLED=0

BINARY_NAME=dbaker
COMPILER_FLAGS_RELEASE=-ldflags "-s -w"

.PHONY: format
format:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: run
run:
	go run ./cmd/${BINARY_NAME}

.PHONY: build
build:
	go build \
		-o ./bin/ \
		./cmd/${BINARY_NAME}


.PHONY: build-release
build-release:
	go build \
		-o ./bin/ \
		${COMPILER_FLAGS_RELEASE} \
		./cmd/${BINARY_NAME}

.PHONY: clean
clean:
	go clean
	rm -rf ./bin/

.PHONY: start-pg
start-pg:
	docker-compose -f test/postgres.docker-compose.yml up
