.PHONY: build vendor_get clean

GOPATH := ${PWD}/_vendor:${GOPATH}

export GOPATH

export http_proxy=http://web-proxy.sgp.hpecorp.net:8080
export https_proxy=https://web-proxy.sgp.hpecorp.net:8080

default: build

build: clean vendor_get main

main:
	go build -v -o suitectl main.go

vendor_get:
	GOPATH=${PWD}/_vendor go get -d -u -v github.com/spf13/cobra/cobra

clean:
	rm -f suitectl 
