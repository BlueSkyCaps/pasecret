package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/config"
)

func main() {
	a := app.New()
	t := &config.DefaultChineseFontTheme{}
	t.SetFonts("STXINWEI.TTF", resourceSTXINWEITTF.StaticContent)
	a.Settings().SetTheme(t)
	w := a.NewWindow("Pasecret")

	w.SetContent(widget.NewLabel("pasecret 世界!"))
	w.ShowAndRun()
}
