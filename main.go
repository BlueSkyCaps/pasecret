package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"pasecret/core/common"
	"pasecret/core/config"
	"pasecret/core/pi18n"
	"pasecret/core/preferences"
	"pasecret/core/storagedata"
	"pasecret/core/ui"
)

func init() {
	storagedata.AppRef.A = app.NewWithID("top.reminisce.pasecret")
	t := &config.DefaultGlobalSettingTheme{}
	t.SetFonts("WenQuanWeiMiHei.ttf", resourceFontWenQuanWeiMiHeiTtf.StaticContent)
	// 更新主题，让fyne使用自定义主题配置
	storagedata.AppRef.A.Settings().SetTheme(t)
	storagedata.AppRef.W = storagedata.AppRef.A.NewWindow("Pasecret")
	storagedata.AppRef.W.CenterOnScreen()
	storagedata.AppRef.W.SetMaster()
	if !fyne.CurrentDevice().IsMobile() {
		//窗体宽度会由子容器grid自动适应
		storagedata.AppRef.W.Resize(fyne.Size{Height: 540})
	}
	//
}

func main() {
	// resource data was bundled by fyne, and it's too long bytes so goland cant analyze immediately
	pi18n.Local12Init(resourceAssetsI18nPasecretZhToml, resourceAssetsI18nPasecretEnToml)
	storagedata.LoadInit(resourceDJson.StaticContent)
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
