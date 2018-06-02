PRJNAME=pw2go
TARGET=bin/$(PRJNAME)
SOURCES=$(wildcard src/$(PRJNAME)/*.go)

.PHONY: run all clean init

all: $(TARGET)
$(TARGET): $(SOURCES)
	gb build

init:
	go get github.com/constabulary/gb/...
	gb vendor restore

run: all
	bin/$(PRJNAME)

clean:
	rm -Rf bin pkg