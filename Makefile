GIT_HASH := $(shell git rev-parse --short HEAD)

dep:
	go mod download

run:
	go run main.go

lint: lint/install lint/run

lint/install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.5.0

lint/run:
	bin/golangci-lint run --config .golangci.yml

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate -f internal/store/sqlc.yaml

build:
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -o supanova-file-cleaner .

docker/local-build:
	DOCKER_BUILDKIT=1 docker build -t supanova-file-cleaner:local .

docker/ci-build:
	DOCKER_BUILDKIT=1 docker build \
	-t supanova-file-cleaner:latest \
	-t supanova-file-cleaner:$(GIT_HASH) .
