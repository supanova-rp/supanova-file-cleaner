run:
	go run main.go

lint: lint-install lint-run

lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.5.0

lint-run:
	bin/golangci-lint run --config .golangci.yml

# build:

# run-docker: