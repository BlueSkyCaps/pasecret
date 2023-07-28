package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"pasecret/core/config"
	"pasecret/core/storagejson"
	"pasecret/core/ui"
)

func init() {
	storagejson.AppRef.A = app.NewWithID("top.reminisce.pasecret")
	t := &config.DefaultGlobalSettingTheme{}
	t.SetFonts("STXINWEI.TTF", resourceSTXINWEITTF.StaticContent)
	// 更新主题，让fyne使用自定义主题配置
	storagejson.AppRef.A.Settings().SetTheme(t)
	storagejson.AppRef.W = storagejson.AppRef.A.NewWindow("Pasecret")
	storagejson.AppRef.W.SetMaster()
	if !fyne.CurrentDevice().IsMobile() {
		storagejson.AppRef.W.Resize(fyne.Size{Height: 500, Width: 650})
	}
}

func main() {
	storagejson.LoadInit(resourceDJson.StaticContent)
	ui.Run()
}
