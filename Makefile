.PHONY: all build-windows build-linux build-debian

all: build

build:
	go build -ldflags="-s -w" -o markdog

build-windows:
	go build -ldflags="-s -w" -o markdog.exe

build-linux:
	go build -ldflags="-s -w" -o markdog

build-debian:
	go build -ldflags="-s -w" -o markdog.deb

run:
	go run main.go
