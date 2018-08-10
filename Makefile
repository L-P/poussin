all: poussin

poussin:
	go build

.PHONY: poussin run pprof
run: poussin
	./poussin boot.gb rom.gb

pprof: poussin
	./poussin -cpuprofile cpu.pprof -memprofile mem.pprof boot.gb rom.gb
	go-torch --file="/tmp/torch.svg" -b cpu.pprof
	sensible-browser "/tmp/torch.svg"
