package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/mattn/go-isatty"
	"github.com/microcosm-cc/bluemonday"
	. "github.com/piax93/pw2go"
	"github.com/zserge/webview"
	"net/url"
	"os"
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
	service = bluemonday.StrictPolicy().Sanitize(service)
	if err := m.pm.AddPassword(service, password, master); err != nil {
		if msg := err.Error(); m.pm.NonFatalError(&msg) {
			(*m.wv).Eval(fmt.Sprintf("modal('%s')", msg))
		} else {
			panic(err)
		}
	} else {
		m.UpdateList()
		(*m.wv).Eval(fmt.Sprintf("modal('Service <i><b>%s</b></i> added')", service))
	}
}

// Get retrieves password from database (webview interface)
func (m *UIManager) Get(service, master string) {
	res, err := m.pm.GetPassword(service, master)
	if err != nil {
		if msg := err.Error(); m.pm.NonFatalError(&msg) {
			(*m.wv).Eval(fmt.Sprintf("modal('%s')", msg))
		} else {
			panic(err)
		}
	} else {
		if err := clipboard.WriteAll(res); err != nil {
			panic(err)
		}
		(*m.wv).Eval("modal('Password copied to clipboard')")
	}
}

// Delete removes password from database (webview interface)
func (m *UIManager) Delete(service string) {
	if err := m.pm.RemovePassword(service); err != nil {
		panic(err)
	}
	m.UpdateList()
}

// SetMaster sets the master password (webview interface)
func (m *UIManager) SetMaster(master string) {
	if err := m.pm.SetMaster(master); err != nil {
		panic(err)
	}
}

// ChangeMaster allows to set a new master password (webview interface)
func (m *UIManager) ChangeMaster(master, newmaster string) {
	if err := m.pm.ChangeMaster(master, newmaster); err != nil {
		if msg := err.Error(); m.pm.NonFatalError(&msg) {
			(*m.wv).Eval(fmt.Sprintf("modal('%s')", msg))
		} else {
			panic(err)
		}
	}
}

// UpdateList updates the service list in the webview
func (m *UIManager) UpdateList() {
	(*m.wv).Eval("clearList()")
	sanitizer := bluemonday.StrictPolicy()
	m.pm.MapServices(func(service string) {
		(*m.wv).Eval(fmt.Sprintf("addService('%s')", sanitizer.Sanitize(service)))
	})
}

// Terminate app
func (m *UIManager) Die() {
	(*m.wv).Terminate()
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

// Start the webview GUI
func startWebUI(pm *PasswordManager) {
	assets, err := loadAssets("app/ui/min/")
	if err != nil {
		panic(err)
	}
	wv := webview.New(webview.Settings{
		URL:       "data:text/html," + url.PathEscape(assets["main.html"]),
		Title:     "Password2Go",
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
		if !pm.IsMasterSet() {
			wv.Eval("setMaster()")
		}
		manager.UpdateList()
	})
	wv.Run()
}

// Start UI
func startUI(pm *PasswordManager, graphical bool) {
	fd := os.Stdout.Fd()
	if (isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)) && !graphical {
		startCLI(pm)
	} else {
		startWebUI(pm)
	}
}
