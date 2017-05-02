NAME = main

all: clean build

build:
	go build ${NAME}.go

clean:
	rm -f ${NAME}

.PHONY: build clean
