package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/storagejsondata"
)

// ShowDataList 点击了某个归类文件夹列表按钮，显示此文件夹的密码项
func ShowDataList(ci storagejson.Category) {
	// 获取此文件夹的密码项
	relatedData := storagejson.GetRelatedDataByCid(ci)
	dataW := storagejson.AppRef.A.NewWindow(ci.Name)
	bottomHBox := container.NewHBox()
	bottomHBox.Add(widget.NewToolbarSpacer().ToolbarObject())
	bottomHBox.Add(widget.NewButton("返回", func() {
		dataW.Close()
	}))
	topHBox := container.NewHBox()
	topHBox.Add(widget.NewToolbarSpacer().ToolbarObject())
	topHBox.Add(widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		showDataEditWin(nil)

	}))

	list := widget.NewList(
		func() int {
			return len(*relatedData)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Button).SetText((*relatedData)[i].Name)
			o.(*widget.Button).OnTapped = func() {
				showDataEditWin(&(*relatedData)[i])
				dialog.ShowInformation("", (*relatedData)[i].Site, storagejson.AppRef.W)
			}
		})
	content := container.NewBorder(topHBox, bottomHBox, nil, nil, list)
	dataW.SetContent(content)
	if !fyne.CurrentDevice().IsMobile() {
		dataW.Resize(fyne.Size{Width: 200, Height: 350})
	}
	dataW.CenterOnScreen()
	dataW.Show()
}
