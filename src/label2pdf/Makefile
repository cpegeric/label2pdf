MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
GOPATH := $(abspath $(MKFILE_PATH)/../../..)
#REV := $(shell git rev-parse --short HEAD)

all:
#	GOPATH=$(GOPATH) go mod init
#	GOPATH=$(GOPATH) go mod tidy
	GOPATH=$(GOPATH) go get .
	GOPATH=$(GOPATH) go build .


clean:
	GOPATH=$(GOPATH) go clean .


.PHONY: all install clean
