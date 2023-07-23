package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"pasecret/core/config"
	"pasecret/core/storagejsondata"
	"pasecret/core/ui"
)

func init() {
	storagejsondata.AppRef.A = app.NewWithID("top.reminisce.test")
	t := &config.DefaultGlobalSettingTheme{}
	t.SetFonts("STXINWEI.TTF", resourceSTXINWEITTF.StaticContent)
	// 更新主题，让fyne使用自定义主题配置
	storagejsondata.AppRef.A.Settings().SetTheme(t)
	storagejsondata.AppRef.W = storagejsondata.AppRef.A.NewWindow("Pasecret")
	storagejsondata.AppRef.W.SetMaster()
	if !fyne.CurrentDevice().IsMobile() {
		storagejsondata.AppRef.W.Resize(fyne.Size{Height: 500, Width: 650})
	}
}

func main() {
	storagejsondata.LoadInit(resourceDJson.StaticContent)
	ui.Run()
}
