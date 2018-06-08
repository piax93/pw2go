package main

import (
	"fmt"
	"github.com/zserge/webview"
	"html/template"
	"net/url"
)

func startUI(pm *PasswordManager) {
	mainui, _ := Asset("ui/min/main.html")
	maincss, _ := Asset("ui/min/main.css")
	bulmacss, _ := Asset("ui/min/bulma.min.css")
	mainjs, _ := Asset("ui/min/main.js")
	wv := webview.New(webview.Settings{
		URL:       "data:text/html," + url.PathEscape(string(mainui)),
		Title:     "PW2GO",
		Width:     400,
		Height:    400,
		Resizable: true,
	})
	wv.Eval(fmt.Sprintf("injectJS(\"%s\");", template.JSEscapeString(string(mainjs))))
	wv.InjectCSS(string(bulmacss))
	wv.InjectCSS(string(maincss))
	for service, _ := range pm.services {
		wv.Eval(fmt.Sprintf("addService('%s')", service))
	}
	wv.Run()
}
