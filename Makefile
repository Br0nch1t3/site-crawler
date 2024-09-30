BIN=./crawler

all: build

build:
	go build .
re: clean build
clean:
	rm -f ${BIN}
