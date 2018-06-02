PRJNAME=pw2go
TARGET=bin/$(PRJNAME)$(EXT)
SOURCES=$(wildcard *.go)

.PHONY: run linux windows deps clean init

linux:
	make $(TARGET) CC=gcc CXX=g++ GOOS=linux

windows:
	make $(TARGET).exe CC=x86_64-w64-mingw32-gcc-win32 CXX=x86_64-w64-mingw32-g++-win32 GOOS=windows EXT=.exe

$(TARGET): $(SOURCES)
	CGO_ENABLED=1 CC=$(CC) CXX=$(CXX) GOOS=$(GOOS) go build
	mkdir -p $(dir $(TARGET))
	mv $(notdir $(TARGET)) $(TARGET)

init:
	go get -u github.com/FiloSottile/gvt

deps: init
	gvt restore
	cd vendor/github.com/andlabs/ui; \
		wget https://github.com/andlabs/ui/raw/master/libui_windows_amd64.a; \
		wget https://github.com/andlabs/ui/raw/master/libui_windows_amd64.res.o; \
		wget https://github.com/andlabs/ui/raw/master/libui_linux_amd64.a;

run:
	bin/$(PRJNAME)$(EXT)

clean:
	rm -Rf bin pkg
