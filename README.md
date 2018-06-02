# Password-2-Go

A simple Go password manager.

Writing this as an exercise to learn the Go programming language.


## Build

#### Dependencies
- **Linux**: `libgtk-3-dev`
- **Windows**: `mingw-w64` toolchain (version 5.0.x)

#### Steps
```bash
cd $GOPATH/src/prjpath   # Enter the project directory
make init                # If you don't have gvt already
make deps                # Download dependencies
make [windows|linux]     # Compile for the desired operating system,
                         # you may need to tweak a little bit the Makefile
                         # depending on your toolchain.
                         # If you are on a Mac, you buy a PC.
```
