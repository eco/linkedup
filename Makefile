GINKGO=ginkgo
LINT=golangci-lint
TEST_PATHS=./...

.DEFAULT_GOAL := default
.PHONY: test init lint test-unit clean

all: bin/lyd bin/lycli bin/ks

test: lint test-unit

lint:
	$(LINT) run
	@echo "Linter passed"

test-unit:
	@echo "Running tests with LCD chain"
	$(GINKGO) $(TEST_PATHS)

init:
	$(MAKE) -C scripts

clean:
	rm -rf bin/

bin/%: cmd/%/*
	go build -o $@ ./$<

default: all init
