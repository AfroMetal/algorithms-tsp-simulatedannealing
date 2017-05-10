NAME = main

all: clean packages build

packages:
	go get github.com/gonum/floats

build:
	go build ${NAME}.go

clean:
	rm -f ${NAME}

.PHONY: build clean
