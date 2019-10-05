GINKGO=ginkgo
LINT=golangci-lint
TEST_PATHS=./...

.PHONY: all test init lint test-unit test-mockchain

install: lyd lycli

test: lint test-mockchain test-unit

lint:
	@$(LINT) run
	@echo "Linter passed"

test-mockchain:
	@echo "Running tests with mock chain"
	@ USE_MOCK_CHAIN=true $(GINKGO) $(TEST_PATHS)

test-unit:
	@echo "Running tests with LCD chain"
	@ USE_MOCK_CHAIN=false $(GINKGO) $(TEST_PATHS)

init:
	$(MAKE) -C scripts
lyd:
	go install ./cmd/lyd

lycli:
	go install ./cmd/lycli

all: install init
