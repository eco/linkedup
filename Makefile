install: lyd lycli

init:
	$(MAKE) -C scripts
lyd:
	go install ./cmd/lyd

lycli:
	go install ./cmd/lycli

all: install init