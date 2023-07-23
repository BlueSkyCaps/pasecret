package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"pasecret/core/config"
	"pasecret/core/storagejsondata"
	"pasecret/core/ui"
)

var appRef storagejsondata.AppStructRef

func init() {
	appRef.A = app.NewWithID("top.reminisce.test")
	t := &config.DefaultGlobalSettingTheme{}
	t.SetFonts("STXINWEI.TTF", resourceSTXINWEITTF.StaticContent)
	// 更新主题，让fyne使用自定义主题配置
	appRef.A.Settings().SetTheme(t)
	appRef.W = appRef.A.NewWindow("Pasecret")
	appRef.W.SetMaster()
	if !fyne.CurrentDevice().IsMobile() {
		appRef.W.Resize(fyne.Size{Height: 500, Width: 650})
	}
}

func main() {
	loadedItems := storagejsondata.LoadInit(appRef, resourceDJson.StaticContent)
	appRef.LoadedItems = loadedItems
	ui.Run(appRef)
}
