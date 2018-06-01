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
	defer pm.Close()
	if len(pm.masterhash) == 0 {
		pm.SetMaster("mastertest")
	}
	fmt.Printf("%v\n", pm)
	err := pm.ChangeMaster("mastertest", "mastertest2")
	fmt.Printf("%v\n", pm)
	if err != nil {
		println(err.Error())
	}
}
