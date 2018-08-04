all: poussin

poussin:
	go build

.PHONY: poussin run
run: poussin
	./poussin
