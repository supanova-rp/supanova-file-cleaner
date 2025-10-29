GIT_HASH := $(shell git rev-parse --short HEAD)
DOCKER_USER := jdgarner
IMAGE_NAME := supanova-file-cleaner

dep:
	go mod download

run:
	go run main.go

lint: lint/install lint/run

lint/install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.5.0

lint/run:
	bin/golangci-lint run --config .golangci.yml

lint/fix:
	bin/golangci-lint run --config .golangci.yml --fix

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate -f internal/store/sqlc.yaml

build/mac:
	CGO_ENABLED=0 \
	GOOS=darwin \
	GOARCH=arm64 \
	go build -o $(IMAGE_NAME) .

build/linux:
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -o $(IMAGE_NAME) .

docker/local-build:
	DOCKER_BUILDKIT=1 docker build -t $(DOCKER_USER)/$(IMAGE_NAME):local .

docker/ci-build:
	DOCKER_BUILDKIT=1 docker build \
	-t $(DOCKER_USER)/$(IMAGE_NAME):latest \
	-t $(DOCKER_USER)/$(IMAGE_NAME):$(GIT_HASH) .

docker/local-run:
	docker run --env-file .env.docker $(DOCKER_USER)/$(IMAGE_NAME):local

docker/push:
	docker push --all-tags $(DOCKER_USER)/$(IMAGE_NAME)
