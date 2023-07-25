package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	lbn := "关闭"
	if fyne.CurrentDevice().IsMobile() {
		lbn = "返回"
	}
	bottomHBox.Add(widget.NewButton(lbn, func() {
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
			return container.NewBorder(
				nil, nil, widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {}), nil,
				widget.NewButton("", func() {}))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			//Container中位置为0的元素是密码项按钮
			(o.(*fyne.Container).Objects[0]).(*widget.Button).SetText((*relatedData)[i].Name)
			(o.(*fyne.Container).Objects[0]).(*widget.Button).OnTapped = func() {
				showDataEditWin(&(*relatedData)[i], ci.Id)
			}
			//Container中位置为1的元素是删除按钮
			(o.(*fyne.Container).Objects[1]).(*widget.Button).OnTapped = func() {
				dialog.ShowConfirm("提示", "确定要删除此条记录吗？", func(b bool) {
					if b {
						storagejson.DeleteData((*relatedData)[i])
					}
				}, storagejson.AppRef.W)
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
		dataW.Resize(fyne.Size{Width: 400, Height: 450})
	}
	dataW.CenterOnScreen()
	dataW.Show()
}
