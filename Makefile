
.PHONY: run all clean

all:
	gb build

run: all
	bin/pw2go

clean:
	rm -Rf bin pkg
