package main

import (
	"fmt"
	"github.com/zserge/webview"
	"net/url"
	"path"
	"strings"
)

// UIManager acts as an interface for receving input from the webview
type UIManager struct {
	pm *PasswordManager
	wv *webview.WebView
}

// Add inserts new password in database (webview interface)
func (m *UIManager) Add(service, password, master string) {

}

// Get retrieves password from database (webview interface)
func (m *UIManager) Get(service, master string) {

}

// Delete removes password from database (webview interface)
func (m *UIManager) Delete(service string) {
	println(service)
}

// SetMaster sets the master password (webview interface)
func (m *UIManager) SetMaster(master string) {

}

// UpdateList updates the service list in the webview
func (m *UIManager) UpdateList() {
	(*m.wv).Eval("clearList()")
	for service := range m.pm.services {
		(*m.wv).Eval(fmt.Sprintf("addService('%s')", service))
	}
}

// Load webview asset file as string in a map
func loadAssets(directory string) (map[string]string, error) {
	res := make(map[string]string)
	files := AssetNames()
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}
	for _, f := range files {
		if strings.HasPrefix(f, directory) {
			content, err := Asset(f)
			if err != nil {
				return res, err
			}
			res[path.Base(f)] = string(content)
		}
	}
	return res, nil
}

// Start the GUI
func startUI(pm *PasswordManager) {
	assets, err := loadAssets("ui/min/")
	if err != nil {
		panic(err)
	}
	wv := webview.New(webview.Settings{
		URL:       "data:text/html," + url.PathEscape(assets["main.html"]),
		Title:     "PW2GO",
		Width:     400,
		Height:    400,
		Resizable: true,
	})
	defer wv.Exit()
	wv.Dispatch(func() {
		manager := UIManager{pm, &wv}
		wv.Bind("manager", &manager)
		wv.InjectCSS(assets["bulma.min.css"])
		wv.InjectCSS(assets["main.css"])
		wv.Eval(assets["utils.js"])
		wv.Eval(assets["main.js"])
		manager.UpdateList()
	})
	wv.Run()
}
