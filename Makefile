PRJNAME=pw2go
TARGET=bin/$(PRJNAME)$(EXT)
SOURCES=$(wildcard *.go)
WINFLAGS="-ldflags='-H windowsgui'"

.PHONY: run linux windows xwindows clean init

linux:
	make $(TARGET) CC=gcc CXX=g++ GOOS=linux

xwindows:
	make $(TARGET).exe CC=x86_64-w64-mingw32-gcc-win32 CXX=x86_64-w64-mingw32-g++-win32 GOOS=windows EXT=.exe FLAGS=$(WINFLAGS)

windows:
	make $(TARGET).exe CC=gcc CXX=g++ GOOS=windows EXT=.exe FLAGS=$(WINFLAGS)

$(TARGET): $(SOURCES)
	go-bindata ui
	CGO_ENABLED=1 CC=$(CC) CXX=$(CXX) GOOS=$(GOOS) go build $(FLAGS)
	mkdir -p $(dir $(TARGET))
	mv $(notdir $(TARGET)) $(TARGET)

init:
	go get -u github.com/FiloSottile/gvt
	go get -u github.com/jteeuwen/go-bindata/...
	gvt restore

run:
	$(TARGET)

clean:
	rm -Rf bin pkg
