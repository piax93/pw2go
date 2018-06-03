package main

import (
	"net/url"

	"github.com/zserge/webview"
)

func startUI(pm *PasswordManager) {
	mainui, _ := Asset("ui/main.html")
	wv := webview.New(webview.Settings{
		URL:       "data:text/html," + url.PathEscape(string(mainui)),
		Title:     "PW2GO",
		Width:     400,
		Height:    400,
		Resizable: true,
	})
	wv.Run()
}
