package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/pi18n"
	"pasecret/core/storagedata"
)

// ShowDataList 点击了某个归类文件夹列表按钮，显示此文件夹的密码项
func ShowDataList(ci storagedata.Category) {
	// 获取当前归类夹 因为如果编辑了归类夹，参数ci的值不是最新的
	realCi := storagedata.GetCategoryByCid(ci.Id)
	// 获取此文件夹的密码项
	relatedData := storagedata.GetRelatedDataByCid(ci.Id)
	dataW := storagedata.AppRef.A.NewWindow(realCi.Name)
	// 将密码项列表窗口贮存到AppRef中后续更新小部件的回调函数使用弹窗作为其父窗口
	storagedata.AppRef.DataListWin = dataW
	bottomHBox := container.NewHBox()
	bottomHBox.Add(widget.NewToolbarSpacer().ToolbarObject())
	lbn := pi18n.LocalizedText("dataListCloseWindowButtonText", nil)
	if fyne.CurrentDevice().IsMobile() {
		lbn = pi18n.LocalizedText("dataListBackWindowButtonText", nil)
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
				dialog.ShowConfirm(pi18n.LocalizedText("dialogShowInformationTitle", nil),
					pi18n.LocalizedText("dataListDeShowConfirm", nil),
					func(b bool) {
						if b {
							storagedata.DeleteData((*relatedData)[i])
						}
					}, dataW)
			}
		})
	// 将定义好的密码项List小部件指向AppRef中贮存
	storagedata.AppRef.DataList = list
	/*
		将ui包中的showDataEditWin方法指向AppRef中贮存，
		后续更新密码项List小部件时需要在storagejson包中更新，storagejson无法循环导入ui包
	*/
	storagedata.AppRef.ShowDataEditWinFunc = showDataEditWin
	content := container.NewBorder(topHBox, bottomHBox, nil, nil, list)
	dataW.SetContent(content)
	if !fyne.CurrentDevice().IsMobile() {
		dataW.Resize(fyne.Size{Width: 400, Height: 450})
	}
	dataW.CenterOnScreen()
	dataW.Show()
}
