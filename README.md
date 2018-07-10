# Password-2-Go

A simple Go password manager.

Writing this as an exercise to learn the Go programming language.


## Build

#### Dependencies
- **Linux**: `gtk-webkit2`
- **Windows**: `mingw-w64` toolchain

#### Steps
```bash
cd $GOPATH/src/prjpath   # Enter the project directory
make init                # Download build tools
make deps                # Download dependencies
make [windows|linux]     # Compile for the desired operating system,
                         # you may need to tweak a little bit the Makefile
                         # depending on your toolchain.
                         # If you are on a Mac, you buy a PC.
```
