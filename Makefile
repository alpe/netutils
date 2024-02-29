.PHONY: test
test: test-unit

.PHONY: test-all
test-all: test-race test-integration

.PHONY: test-unit
test-unit:
	go test -mod=readonly ./...

.PHONY: test-race
test-race:
	go test -mod=readonly -race ./...

.PHONY: lint
lint:
	golangci-lint run -v ./... --timeout 5m
