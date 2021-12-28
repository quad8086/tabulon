GO=go
VERSION=dev
SRCDIR=$(shell pwd)
INSTALL_ROOT=$(SRCDIR)
INSTALL_DIR=$(INSTALL_ROOT)/.install
TARGET=tabulon

.NOTPARALLEL:

all: deps build

deps:
	$(GO) mod download github.com/jessevdk/go-flags
	$(GO) mod download github.com/gdamore/tcell
	$(GO) mod download github.com/gdamore/tcell/v2
	$(GO) mod download github.com/danielgtaylor/mexpr
	$(GO) get github.com/jessevdk/go-flags
	$(GO) get github.com/gdamore/tcell/v2

build:
	$(GO) build

clean:
	$(RM) $(TARGET) go.sum

install:
	mkdir -p $(INSTALL_DIR)/bin
	install --mode 755 $(TARGET) $(INSTALL_DIR)/bin
