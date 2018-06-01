package main

import (
	"fmt"
	"log"
)

func printerr(err error) {
	if err != nil {
		println("Error: " + err.Error())
	}
}

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
	printerr(pm.ChangeMaster("mastertest", "mastertest2"))
	fmt.Printf("%v\n", pm)
	printerr(pm.AddPassword("testservice", "password123", "mastertest2"))
	plain, err := pm.GetPassword("testservice", "mastertest2")
	printerr(err)
	println("PASS: " + plain)
}
