package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"pasecret/core/config"
	"pasecret/core/ui"
)

var w fyne.Window
var a fyne.App

func init() {
	a = app.NewWithID("top.reminisce.test")
	t := &config.DefaultChineseFontTheme{}
	t.SetFonts("STXINWEI.TTF", resourceSTXINWEITTF.StaticContent)
	a.Settings().SetTheme(t)
	w = a.NewWindow("Pasecret")
	w.SetMaster()
}
func main() {
	ui.Run(w, a)
}
