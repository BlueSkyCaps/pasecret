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
	t := &config.DefaultGlobalSettingTheme{}
	t.SetFonts("STXINWEI.TTF", resourceSTXINWEITTF.StaticContent)
	// 更新主题，让fyne使用自定义主题配置
	a.Settings().SetTheme(t)
	w = a.NewWindow("Pasecret")
	w.SetMaster()
}
func main() {
	ui.Run(w, a)
}
