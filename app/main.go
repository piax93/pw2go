package main

import (
	"flag"
	. "github.com/piax93/pw2go"
)

// Main function
func main() {
	println("---- PW2GO, a Go password manager ----")
	graphical := flag.Bool("g", false, "Force to use webview GUI")
	flag.Parse()
	pm := NewPassManager("pw2go", "master", "passwords")
	if err := pm.Init(); err != nil {
		panic(err)
	}
	defer pm.Close()
	startUI(pm, *graphical)
}
