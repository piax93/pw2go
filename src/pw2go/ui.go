package main

import (
	"github.com/andlabs/ui"
)

func startUI(pm *PasswordManager) {
	if err := ui.Main(func() {

	}); err != nil {
		panic(err)
	}
}
