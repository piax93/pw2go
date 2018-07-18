PRJNAME=pw2go
TARGET=app/bin/$(PRJNAME)$(EXT)
SOURCES=$(wildcard *.go) $(wildcard app/ui/*) $(wildcard app/*.go) Makefile
LINUXFLAGS="-ldflags='-w -s'"
WINFLAGS="-ldflags='-w -s -H windowsgui'"

.PHONY: run linux windows xwindows clean init

linux:
	make $(TARGET) CC=gcc CXX=g++ GOOS=linux FLAGS=$(LINUXFLAGS)
	strip -s $(TARGET)

xwindows:
	make $(TARGET).exe CC=x86_64-w64-mingw32-gcc-win32 CXX=x86_64-w64-mingw32-g++-win32 GOOS=windows EXT=.exe FLAGS=$(WINFLAGS)

windows:
	make $(TARGET).exe CC=gcc CXX=g++ GOOS=windows EXT=.exe FLAGS=$(WINFLAGS)

$(TARGET): $(SOURCES)
	mkdir -p app/ui/min
	minify -o app/ui/min app/ui
	go-bindata -o app/bindata.go app/ui/min
	mkdir -p $(dir $(TARGET))
	CGO_ENABLED=1 CC=$(CC) CXX=$(CXX) GOOS=$(GOOS) go build $(FLAGS) -o $(TARGET) ./app

init:
	go get -u github.com/FiloSottile/gvt
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/tdewolff/minify/cmd/minify
	gvt restore

run:
	$(TARGET)

clean:
	rm -Rf app/bin
