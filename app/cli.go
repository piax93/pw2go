package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	. "github.com/piax93/pw2go"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"regexp"
)

// Print options
func printMenu() {
	println("\n")
	println("1) Add password")
	println("2) Get password")
	println("3) List services")
	println("4) Change master")
	println("5) Exit")
}

// Read input
func getInput(prompt string, confidential bool) string {
	var res string
	print(prompt)
	if confidential {
		bytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		res = string(bytes)
		print("\n")
	} else {
		if _, err := fmt.Scanln(&res); err != nil {
			panic(err)
		}
	}
	return res
}

// Print error message is non-critical, otherwise panic
func printOrPanic(pm *PasswordManager, err *error) {
	if msg := (*err).Error(); pm.NonFatalError(&msg) {
		println(msg)
	} else {
		panic(err)
	}
}

// Start command line interface
func startCLI(pm *PasswordManager) {
	if !pm.IsMasterSet() {
		pm.SetMaster(getInput("Set master password: ", true))
	}
	running := true
	for running {
		printMenu()
		choice := getInput("Choice: ", false)
		if res, err := regexp.MatchString("[1-5]", choice); err != nil || !res {
			println("Bad input")
			continue
		}
		switch choice {
		case "1":
			if err := pm.AddPassword(
				getInput("Service name: ", false),
				getInput("Service password: ", true),
				getInput("Master: ", true),
			); err != nil {
				printOrPanic(pm, &err)
			} else {
				println("Password corretly stored")
			}
		case "2":
			if res, err := pm.GetPassword(
				getInput("Service name: ", false),
				getInput("Master: ", true),
			); err != nil {
				printOrPanic(pm, &err)
			} else {
				if err := clipboard.WriteAll(res); err != nil {
					panic(err)
				}
				println("Password copied to clipboard")
			}
		case "3":
			println("Services:")
			pm.MapServices(func(service string) {
				fmt.Printf("- %s\n", service)
			})
		case "4":
			if err := pm.ChangeMaster(
				getInput("Old master password: ", true),
				getInput("New master password: ", true),
			); err != nil {
				printOrPanic(pm, &err)
			} else {
				println("Master password changed")
			}
		case "5":
			running = false
			println("Bye bye")
		}

	}
}
