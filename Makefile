VERSION = $(shell git describe --always --long --dirty)
PACKAGE = github.com/n0ot/nextpass
BUILDFLAGS = -ldflags "-X $(PACKAGE).Version=$(VERSION)"
.PHONY: clean install nextpass

nextpass:
	go build $(BUILDFLAGS) -o nextpass cmd/nextpass/main.go

install:
	go install $(BUILDFLAGS) $(PACKAGE)/cmd/nextpass

clean:
	rm -f nextpass