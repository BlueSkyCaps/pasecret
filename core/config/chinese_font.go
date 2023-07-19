package config

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"golang.org/x/image/colornames"
	"image/color"
)

// DefaultGlobalSettingTheme 全局配置支持中文字体，避免乱码，以及一些个性化主题颜色等
type DefaultGlobalSettingTheme struct {
	regular, bold, italic, boldItalic, monospace fyne.Resource
}

// Color Theme接口必须实现的方法
func (t *DefaultGlobalSettingTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// 设置全局主题主调色-淡绿
	if name == theme.ColorNamePrimary {
		return color.RGBA{R: 0, G: 128, B: 0, A: 255}
	}
	// 设置全局鼠标滑入显示的颜色
	if name == theme.ColorNameHover {
		return colornames.Aliceblue
	}
	// 设置全局点击时显示的颜色
	if name == theme.ColorNamePressed {
		return colornames.Lightblue
	}
	return theme.DefaultTheme().Color(name, variant)
}

// Icon Theme接口必须实现的方法
func (t *DefaultGlobalSettingTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Font Theme接口必须实现的方法
func (m *DefaultGlobalSettingTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return m.monospace
	}
	if style.Bold {
		if style.Italic {
			return m.boldItalic
		}
		return m.bold
	}
	if style.Italic {
		return m.italic
	}
	return m.regular
}

// Size Theme接口必须实现的方法
func (m *DefaultGlobalSettingTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// SetFonts 主动更改默认字体，若不调用次函数，字体不会更改。
func (t *DefaultGlobalSettingTheme) SetFonts(staticName string, staticContent []byte) {
	t.regular = theme.TextFont()
	t.bold = theme.TextBoldFont()
	t.italic = theme.TextItalicFont()
	t.boldItalic = theme.TextBoldItalicFont()
	t.monospace = theme.TextMonospaceFont()

	if staticName != "" {
		t.regular = loadCustomFont(staticName, staticContent)
		t.bold = loadCustomFont(staticName, staticContent)
		t.italic = loadCustomFont(staticName, staticContent)
		t.boldItalic = loadCustomFont(staticName, staticContent)
		t.monospace = t.regular
	}
}

func loadCustomFont(env string, content []byte) fyne.Resource {

	resource := fyne.NewStaticResource(env, content)
	return resource
}
