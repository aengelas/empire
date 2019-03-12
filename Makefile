.PHONY: build test bootstrap

REPO = remind101/empire
TYPE ?= patch
ARTIFACTS ?= build

cmds: build/empire build/emp

clean:
	rm -rf build/*

build/empire:
	# Setting -tags netgo will force the binary to use the Go networking
	# libraries instead of system libraries, for more consistent behavior.
	# The -ldflags setting will enable compiling a static binary that's
	# portable regardless of whether libc is available.
	go build -ldflags '-extldflags "-static"' -tags netgo -o build/empire ./cmd/empire

build/emp:
	# Setting -tags netgo will force the binary to use the Go networking
	# libraries instead of system libraries, for more consistent behavior.
	# The -ldflags setting will enable compiling a static binary that's
	# portable regardless of whether libc is available.
	go build -ldflags '-extldflags "-static"' -tags netgo -o build/emp ./cmd/emp

bootstrap: cmds
	createdb empire || true
	./build/empire migrate

build: Dockerfile
	docker build -t ${REPO} .

test: build/emp
	go test -race $(shell go list ./... | grep -v /vendor/)
	./tests/deps

vet:
	go vet $(shell go list ./... | grep -v /vendor/)

bump:
	pip install --upgrade bumpversion
	bumpversion ${TYPE}

$(ARTIFACTS)/all: $(ARTIFACTS)/emp-Linux-x86_64 $(ARTIFACTS)/emp-Darwin-x86_64 $(ARTIFACTS)/empire-Linux-x86_64

$(ARTIFACTS)/emp-Linux-x86_64:
	env GOOS=linux go build -o $@ ./cmd/emp
$(ARTIFACTS)/emp-Darwin-x86_64:
	env GOOS=darwin go build -o $@ ./cmd/emp

$(ARTIFACTS)/empire-Linux-x86_64:
	env GOOS=linux go build -o $@ ./cmd/empire
