PLATFORM ?= $(shell uname -m)
OSTYPE := $(shell uname | tr "[A-Z]" "[a-z]")
DISTDIR = ./

GOOS ?= $(OSTYPE)
CGO_ENABLED = 0

export PLATFORM
export GOOS
export CGO_ENABLED


build:
	@[ -d $(DISTDIR) ] || mkdir -p $(DISTDIR)
	go build --tags "static" -v -ldflags '-s -w' -o $(DISTDIR)/helloapp-$(GOOS)-$(PLATFORM)

build-%:
	GOOS=$* go build --tags "static" -v -ldflags '-s -w' -o $(DISTDIR)/helloapp-$*-$(PLATFORM)

dist-shasum: DISTDIR=./dist/
dist-shasum:
	cd $(DISTDIR) && sha256sum helloapp-$(GOOS)-$(PLATFORM) | tee helloapp-$(GOOS)-$(PLATFORM).sha256

dist: DISTDIR=./dist/
dist: build

dist-%: DISTDIR=./dist/
dist-%: build-%
	$(NOOP)

lint:
	[ -x ./bin/golangci-lint ] || wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s
	./bin/golangci-lint run

test: lint
	go test
