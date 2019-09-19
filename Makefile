install: ly lycli

init:
	$(MAKE) -C scripts
ly:
	go install ./cmd/ly

lycli:
	go install ./cmd/lycli

all: install init