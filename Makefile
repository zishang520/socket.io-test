.PHONY: default update install all engine.io socket.io init fmt
export GOPATH:=$(shell pwd)/vendor
# Set the GOPROXY environment variable
export GOPROXY=https://goproxy.io,direct
# export http_proxy=socks5://127.0.0.1:1080
# export https_proxy=%http_proxy%

default: all

install:
	go mod tidy -v
	go mod vendor -v

update:
	go mod tidy -v

fmt:
	go fmt ...

engine.io:
	go build --mod=mod -o "bin/engine" main.go
	"bin/main"

socket.io:
	go build --mod=mod -o "bin/socket" main.go
	"bin/main"

init:


all:
	CGO_ENABLED=0 go build --mod=mod  -ldflags '-s -w -extldflags "-static"' -o "bin/engineio" main.go
