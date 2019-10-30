GINKGO=ginkgo
LINT=golangci-lint
TEST_PATHS=./...

.DEFAULT_GOAL := default
.PHONY: test init lint test-unit clean redeploy

all: bin/lyd bin/lycli bin/ks

test: lint test-unit

lint:
	$(LINT) run
	@echo "Linter passed"

test-unit:
	@echo "Running tests with LCD chain"
	$(GINKGO) $(TEST_PATHS)

init: bin/lyd bin/lycli scripts/initChain.sh
	scripts/initChain.sh

clean:
	rm -rf bin/

bin/%: cmd/%/* $(shell find x/ -type f)
	go build -o $@ ./$<

default: all init

redeploy: bin/lyd
	echo "Exporting current state of app"
	@bin/lyd export --for-zero-height > lyd_export.json

	echo "Reseting the chain state to 0"
	@bin/lyd unsafe-reset-all

	echo "Moving exported state to genesis"
	cp lyd_export.json ~/.lyd/config/genesis.json

	echo "Restarting the chain"
	bin/lyd start
