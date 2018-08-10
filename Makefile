all: poussin

poussin:
	go build

.PHONY: poussin run pprof
run: poussin
	./poussin boot.gb rom.gb

pprof: poussin
	./poussin -cpuprofile cpu.pprof boot.gb rom.gb
	go-torch --file="/tmp/cpu.svg" -b cpu.pprof
	sensible-browser "/tmp/cpu.svg"
