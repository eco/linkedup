GINKGO=ginkgo
LINT=golangci-lint
TEST_PATHS=./...

.PHONY: all test init lint test-unit

install: lyd lycli ks

test: lint test-unit

lint:
	$(LINT) run
	@echo "Linter passed"

test-unit:
	@echo "Running tests with LCD chain"
	$(GINKGO) $(TEST_PATHS)

init:
	$(MAKE) -C scripts
lyd:
	go build -o ./bin/lyd ./cmd/lyd

lycli:
	go build -o ./bin/lycli ./cmd/lycli

ks:
	go build -o ./bin/ks ./cmd/ks

all: install init
