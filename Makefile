
BASE_PACKAGE=github.com/wselwood/backuptool
TOOLS=backup

MAKEPATH := $(abspath $(lastword $(MAKEFILE_LIST)))
BASEDIR := $(dir $(MAKEPATH))

GOPATH=GOPATH=$(BASEDIR)
GO=$(GOPATH) go

all: $(TOOLS)

$(TOOLS):%: test
	$(GO) build -o $@ $(BASE_PACKAGE)/cmd/$@

test: 
	$(GO) test -race -cover $(BASE_PACKAGE)/...

dep:
	cd $(BASEDIR)src/$(BASE_PACKAGE); $(GOPATH) dep ensure

clean:
	rm -f $(TOOLS)

code:
	$(GOPATH) code .
