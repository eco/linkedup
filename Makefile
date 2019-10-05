GINKGO=ginkgo
LINT=golangci-lint
TEST_PATHS=./...

.PHONY: all test init lint test-unit

install: lyd lycli

test: lint test-unit

lint:
	$(LINT) run

test-unit:
	USE_MOCK_CHAIN=true $(GINKGO) $(TEST_PATHS)
	$(GINKGO) $(TEST_PATHS)

init:
	$(MAKE) -C scripts
lyd:
	go install ./cmd/lyd

lycli:
	go install ./cmd/lycli

all: install init
