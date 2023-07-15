// Package config DefaultChineseFontTheme 配置支持中文字体 避免乱码
package config

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type DefaultChineseFontTheme struct {
	regular, bold, italic, boldItalic, monospace fyne.Resource
}

func (t *DefaultChineseFontTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t *DefaultChineseFontTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *DefaultChineseFontTheme) Font(style fyne.TextStyle) fyne.Resource {
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

func (m *DefaultChineseFontTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (t *DefaultChineseFontTheme) SetFonts(staticName string, staticContent []byte) {
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
