package main

func main() {
	println("---- PW2GO, a Go password manager ----")
	pm := PasswordManager{dbname: "pw2go", mastertable: "master", passtable: "passwords"}
	if err := pm.Init(); err != nil {
		panic(err)
	}
	defer pm.Close()
	if len(pm.masterhash) == 0 {
		pm.SetMaster("mastertest")
	}
	// startUI(&pm)
}
