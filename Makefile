VERSION = $(shell git describe --always --long --dirty)
PACKAGE = github.com/n0ot/nextpass
BUILDFLAGS = -ldflags "-X $(PACKAGE).Version=$(VERSION)"
ifdef WINDIR
	CURRENT = nextpass.exe
else
	CURRENT = nextpass
endif

.PHONY: all clean install windows linux freebsd darwin

nextpass: current

all: windows linux freebsd darwin current

windows: bin
	GOOS=windows GOARCH=amd64 go build $(BUILDFLAGS) -o bin/nextpass-windows-amd64.exe cmd/nextpass/main.go

linux: bin
	GOOS=linux GOARCH=amd64 go build $(BUILDFLAGS) -o bin/nextpass-linux-amd64 cmd/nextpass/main.go

freebsd: bin
	GOOS=freebsd GOARCH=amd64 go build $(BUILDFLAGS) -o bin/nextpass-freebsd-amd64 cmd/nextpass/main.go

darwin: bin
	GOOS=darwin GOARCH=amd64 go build $(BUILDFLAGS) -o bin/nextpass-darwin-amd64 cmd/nextpass/main.go

current: bin
	go build $(BUILDFLAGS) -o bin/$(CURRENT) cmd/nextpass/main.go

bin:
	mkdir bin

install:
	go install $(BUILDFLAGS) $(PACKAGE)/cmd/nextpass

clean:
	rm -rf bin