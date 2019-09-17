install: ly lycli

init:
	$(MAKE) -C scripts
bt:
	go install ./cmd/bly

btcli:
	go install ./cmd/lycli

all: install init