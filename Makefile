all: clean build

clean:
	rm -rf ./bin

build:
	go build -o ./bin/afl-multicore ./cmd/afl-multicore
	go build -o ./bin/afl-multikill ./cmd/afl-multikill
	go build -o ./bin/afl-multistats ./cmd/afl-multistats

install:
	mv ./bin/* ${GOBIN}