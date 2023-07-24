package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/storagejson"
)

// ShowDataList 点击了某个归类文件夹列表按钮，显示此文件夹的密码项
func ShowDataList(ci storagejson.Category) {
	// 获取此文件夹的密码项
	relatedData := storagejson.GetRelatedDataByCid(ci.Id)
	dataW := storagejson.AppRef.A.NewWindow(ci.Name)
	bottomHBox := container.NewHBox()
	bottomHBox.Add(widget.NewToolbarSpacer().ToolbarObject())
	bottomHBox.Add(widget.NewButton("返回", func() {
		dataW.Close()
	}))
	topHBox := container.NewHBox()
	topHBox.Add(widget.NewToolbarSpacer().ToolbarObject())
	topHBox.Add(widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		showDataEditWin(nil, ci.Id)
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
				showDataEditWin(&(*relatedData)[i], ci.Id)
			}
		})
	// 将定义好的密码项List小部件指向AppRef中贮存
	storagejson.AppRef.DataList = list
	/*
		将ui包中的showDataEditWin方法指向AppRef中贮存，
		后续更新密码项List小部件时需要在storagejson包中更新，storagejson无法循环导入ui包
	*/
	storagejson.AppRef.ShowDataEditWinFunc = showDataEditWin
	content := container.NewBorder(topHBox, bottomHBox, nil, nil, list)
	dataW.SetContent(content)
	if !fyne.CurrentDevice().IsMobile() {
		dataW.Resize(fyne.Size{Width: 200, Height: 350})
	}
	dataW.CenterOnScreen()
	dataW.Show()
}
