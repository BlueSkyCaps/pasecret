package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"pasecret/core/storagejsondata"
)

// ShowCategoryAddWin 首页点击了添加按钮，添加新归类夹
func ShowCategoryAddWin() {
	addW := storagejson.AppRef.A.NewWindow("编辑归类")
	vBox := container.NewVBox()
	vBox.Add(widget.NewLabel("名称："))
	nameEntry := widget.NewEntry()
	vBox.Add(nameEntry)
	vBox.Add(widget.NewLabel("描述："))
	descriptionEntry := widget.NewEntry()
	vBox.Add(descriptionEntry)
	vBox.Add(widget.NewLabel("Alias："))
	aliasEntry := widget.NewEntry()
	vBox.Add(aliasEntry)
	// 创建水平按钮布局
	editCancelBtn := widget.NewButton("取消", func() {
		addW.Close()
	})
	editConfirmBtn := widget.NewButton("确定", func() {
		e := &common.EditForm{Name: nameEntry.Text, Alias: aliasEntry.Text, Description: descriptionEntry.Text}
		AddConfirm(e)
		addW.Close()
	})
	editConfirmBtn.Importance = widget.HighImportance
	hBox := container.NewHBox(widget.NewToolbarSpacer().ToolbarObject(), editCancelBtn, editConfirmBtn)
	vBox.Add(hBox)
	addW.SetContent(vBox)
	if !fyne.CurrentDevice().IsMobile() {
		addW.Resize(fyne.Size{Width: 300, Height: vBox.Size().Height})
	}
	addW.CenterOnScreen()
	addW.Show()
}

func AddConfirm(e *common.EditForm) {
	if common.IsWhiteAndSpace(e.Name) {
		dialog.ShowInformation("提示", "归类文件夹名称不能是空的。", storagejson.AppRef.W)
		return
	}
	// 先成功新增到本地存储库
	addCi := storagejson.AddCategory(e)
	if addCi == nil {
		return
	}
	// 再根据Category对象创建引用Cart小部件
	cart := CreateCurrentCart(*addCi)
	// 最后更新当前Cart小部件布局
	storagejson.AppRef.RepaintCartsByAdd(*addCi, cart)
}
