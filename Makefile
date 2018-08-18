all: poussin

poussin:
	go build

.PHONY: poussin run pprof test
run: poussin
	./poussin rom.gb  2> /dev/shm/stderr

pprof: poussin
	./poussin -cpuprofile cpu.pprof rom.gb 2> /dev/shm/stderr
	go tool pprof -web cpu.pprof

test:
	go vet ./...
	go test . -timeout 15s
