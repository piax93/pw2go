package main

// Main function
func main() {
	println("---- PW2GO, a Go password manager ----")
	pm := PasswordManager{dbname: "pw2go", mastertable: "master", passtable: "passwords"}
	if err := pm.Init(); err != nil {
		panic(err)
	}
	defer pm.Close()
	/*
		if len(pm.masterhash) == 0 {
			pm.SetMaster("mastertest")
		}
		pm.AddPassword("service1", "password1", "mastertest")
		pm.AddPassword("service2", "password2", "mastertest")
	*/
	startUI(&pm)
}
