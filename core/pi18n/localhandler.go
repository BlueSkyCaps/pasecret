package pi18n

// 根据首选项配置加载中文/英文本地化文本

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"pasecret/core/common"
	"pasecret/core/preferences"
	"pasecret/core/storagedata"
	"path"
	"time"
)

var localize *i18n.Localizer
var localFilePath string
var Lang string

// Local12Init 不存在本地化文本文件则创建，并且初始化本地语言环境。
func Local12Init(zhToml *fyne.StaticResource, enToml *fyne.StaticResource) {
	// 获取首选项中的当前语言
	Lang = preferences.GetPreferenceByLocalLang()
	// 如果为空，则是首次安装打开（以及低版本pasecret没有语言设置），默认显示中文
	var localFileBytes []byte
	if common.IsWhiteAndSpace(Lang) || Lang == "zh" {
		localFilePath = path.Join(storagedata.AppRef.A.Storage().RootURI().Path(), "pasecret.zh.toml")
		Lang = "zh"
		// 从bundled.go读取打包的语言文件字节数据
		localFileBytes = zhToml.StaticContent
	} else if Lang == "en" {
		localFilePath = path.Join(storagedata.AppRef.A.Storage().RootURI().Path(), "pasecret.en.toml")
		localFileBytes = enToml.StaticContent
	}
	// 不存在当前语言本地化文件，则创建
	if !common.Existed(localFilePath) {
		r, err := common.CreateFile(localFilePath, localFileBytes)
		if !r {
			dialog.NewInformation("err", "Local12Init, CreateFile:\n"+err.Error(), storagedata.AppRef.W).Show()
			go func() {
				time.Sleep(time.Second * 2)
				os.Exit(1)
			}()
			return
		}
	}
	initLocalize(Lang)
	/*
		将pi18n包中的LocalizedText方法指向AppRef中贮存，
		后续更新密码项List小部件时需要在storagejson包中更新，storagejson无法循环导入pi18n包
	*/
	storagedata.AppRef.LocalizedTextFunc = LocalizedText
}

// 配置当前全局语言环境
func initLocalize(lang string) {
	// 创建Bundle
	bundle := i18n.NewBundle(language.Chinese)
	// 注册解析器
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// 加载当前语言文件
	_, err := bundle.LoadMessageFile(localFilePath)
	if err != nil {
		dialog.NewInformation("err", "initLocalize, LoadMessageFile:\n"+err.Error(), storagedata.AppRef.W).Show()
		go func() {
			time.Sleep(time.Second * 2)
			os.Exit(1)
		}()
		return
	}
	localize = i18n.NewLocalizer(bundle, lang)
}

// LocalizedText 获取本地化文本
func LocalizedText(messageId string, tmd map[string]interface{}) string {
	i := i18n.LocalizeConfig{}
	if tmd != nil {
		i = i18n.LocalizeConfig{
			MessageID:    messageId,
			TemplateData: tmd,
		}
	} else {
		i = i18n.LocalizeConfig{
			MessageID: messageId,
		}
	}
	translation, err := localize.Localize(&i)
	if err != nil {
		dialog.NewInformation("err", "LocalizedText, Localize:\n"+err.Error(), storagedata.AppRef.W).Show()
		go func() {
			time.Sleep(time.Second * 2)
			os.Exit(1)
		}()
	}
	return translation
}
