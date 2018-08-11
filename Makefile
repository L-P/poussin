all: poussin

poussin:
	go build

.PHONY: poussin run pprof test
run: poussin
	./poussin boot.gb rom.gb  2> stderr

pprof: poussin
	./poussin -cpuprofile cpu.pprof boot.gb rom.gb 2> stderr
	go-torch --file="/tmp/cpu.svg" -b cpu.pprof
	sensible-browser "/tmp/cpu.svg"

test:
	go vet ./...
	golint ./...
