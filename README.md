# Password-2-Go

[![Build Status](https://travis-ci.org/piax93/pw2go.svg?branch=master)](https://travis-ci.org/piax93/pw2go)

A simple Go password manager.

Writing this as an exercise to learn the Go programming language.


## Build

#### Dependencies
- **Linux**: `libwebkit2gtk-4.0-dev`
- **Windows**: `mingw-w64` toolchain

#### Steps
```bash
cd $GOPATH/src/prjpath   # Enter the project directory
make init                # Download build tools and dependencies
make [windows|linux]     # Compile for the desired operating system,
                         # you may need to tweak a little bit the Makefile
                         # depending on your toolchain.
                         # If you are on a Mac, you buy a PC.
```
