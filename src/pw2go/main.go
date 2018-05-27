package main

import (
	"fmt"
	"log"
)

func main() {
	println("--- PW2GO, a Go password manager ---")
	pm := PasswordManager{dbname: "pw2go", mastertable: "master", passtable: "passwords"}
	if err := pm.Init(); err != nil {
		log.Fatal(err)
	}
	if len(pm.masterhash) == 0 {
		pm.SetMaster("mastertest")
	}
	defer pm.Close()
	fmt.Printf("%v\n", pm)
}
