// Package ui 定义窗口主要UI内容元素
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Run 开始定义UI元素并显示窗口
func Run(w fyne.Window, a fyne.App) {
	// tabs
	firstAddTabs(w, a)
	w.ShowAndRun()
}

func firstAddTabs(w fyne.Window, a fyne.App) {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.HomeIcon(), widget.NewLabel("Hello")),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), widget.NewButton("Open new", func() {
			w3 := a.NewWindow("Third")
			w3.SetContent(widget.NewButton("确定", func() {
				w3.Close()
			}))

			w3.Show()
		})),
	)
	tabs.SetTabLocation(container.TabLocationBottom)
	w.SetContent(tabs)
}
