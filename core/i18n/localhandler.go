// Package i18n 根据首选项配置加载中文/英文本地化文本
package i18n

//import (
//	"encoding/json"
//	"fyne.io/fyne/v2/dialog"
//	"github.com/BurntSushi/toml"
//	"github.com/nicksnyder/go-i18n/v2/i18n"
//	"golang.org/x/text/language"
//	"pasecret/core/common"
//	"pasecret/core/storagejson"
//	"path"
//)
//
//// Local12Init 不存在本地化文本文件则创建
//func Local12Init() {
//	preferencePath = path.Join(storagejson.AppRef.A.Storage().RootURI().Path(), "preference.json")
//	// 不存在首选项文件，则创建
//	if !common.Existed(preferencePath) {
//		initPre := Preferences{
//			LockPwd: "",
//		}
//		marshal, err := json.Marshal(initPre)
//		if err != nil {
//			dialog.ShowInformation("err", "preferenceInit, json.Marshal:"+err.Error(), storagejson.AppRef.W)
//			return
//		}
//		r, err := common.CreateFile(preferencePath, marshal)
//		if !r {
//			dialog.NewInformation("err", "preferenceInit, CreateFile:"+err.Error(), storagejson.AppRef.W).Show()
//			return
//		}
//	}
//}
//func Test01(messageId string, lang string) string {
//	// 创建 Bundle
//	bundle := i18n.NewBundle(language.English)
//
//	// 注册解析器
//	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
//
//	// 创建 Localize
//	var localize *i18n.Localizer
//	if lang == "en" {
//		_, err := bundle.LoadMessageFile("assets/i18n/pasecret.en.toml")
//		if err != nil {
//			return ""
//		}
//		localize = i18n.NewLocalizer(bundle, "en")
//
//	} else {
//		_, err := bundle.LoadMessageFile("assets/i18n/pasecret.zh.toml")
//		if err != nil {
//			return ""
//		}
//
//		localize = i18n.NewLocalizer(bundle, "zh")
//	}
//
//	// 执行翻译
//	translation, err := localize.Localize(&i18n.LocalizeConfig{
//		MessageID: messageId,
//		TemplateData: map[string]interface{}{
//			"Name":  "Nick",
//			"Count": 2,
//		},
//	})
//	if err != nil {
//		return ""
//	}
//	return translation
//}
