package main

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func Test01(messageId string, lang string) string {
	// 创建 Bundle
	bundle := i18n.NewBundle(language.English)

	// 注册解析器
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 创建 Localize
	var localize *i18n.Localizer
	if lang == "en" {
		_, err := bundle.LoadMessageFile("assets/pi18n/pasecret.en.toml")
		if err != nil {
			return err.Error()
		}
		localize = i18n.NewLocalizer(bundle, lang)

	} else {
		_, err := bundle.LoadMessageFile("assets/pi18n/pasecret.zh.toml")
		if err != nil {
			return err.Error()
		}

		localize = i18n.NewLocalizer(bundle, lang)
	}
	// 执行翻译
	translation, err := localize.Localize(&i18n.LocalizeConfig{
		MessageID: messageId,
	})

	if err != nil {
		return err.Error()
	}
	return translation
}
