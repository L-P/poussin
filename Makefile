all: poussin

poussin:
	go build

.PHONY: poussin run pprof test
run: poussin
	./poussin boot.gb rom.gb  2> stderr

pprof: poussin
	./poussin -cpuprofile cpu.pprof boot.gb rom.gb 2> stderr
	go tool pprof -web cpu.pprof

test:
	go vet ./...
	golint ./...
