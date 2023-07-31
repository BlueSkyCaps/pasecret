package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"pasecret/core/common"
	"pasecret/core/config"
	"pasecret/core/preferences"
	"pasecret/core/storagejson"
	"pasecret/core/ui"
)

func init() {

	storagejson.AppRef.A = app.NewWithID("top.reminisce.pasecret")
	t := &config.DefaultGlobalSettingTheme{}
	t.SetFonts("WenQuanWeiMiHei.ttf", resourceWenQuanWeiMiHeiTtf.StaticContent)
	// 更新主题，让fyne使用自定义主题配置
	storagejson.AppRef.A.Settings().SetTheme(t)
	storagejson.AppRef.W = storagejson.AppRef.A.NewWindow("Pasecret")
	storagejson.AppRef.W.CenterOnScreen()
	storagejson.AppRef.W.SetMaster()
	if !fyne.CurrentDevice().IsMobile() {
		//窗体宽度会由子容器grid自动适应
		storagejson.AppRef.W.Resize(fyne.Size{Height: 540})
	}
}

func main() {
	storagejson.LoadInit(resourceDJson.StaticContent)
	uIHandler()
}

func uIHandler() {
	// 有设置启动密码则先显示解锁
	if !common.IsWhiteAndSpace(preferences.GetPreferenceByLockPwd()) {
		ui.Run(true)
	} else {
		// 否则直接显示主窗口
		ui.Run(false)
	}
}
