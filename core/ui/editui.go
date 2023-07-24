package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"pasecret/core/storagejsondata"
)

// ShowCategoryEditWin 点击了某归类文件夹Card的Edit按钮，显示此具体编辑窗口
func ShowCategoryEditWin(ci storagejson.Category, ciCard *widget.Card) {
	var realCi storagejson.Category
	// 从AppRef中根据id找到文件夹，因为回调函数传的Category参数是最初原始的数据，而AppRef中LoadedItems是实时更新的数据
	for _, nci := range storagejson.AppRef.LoadedItems.Category {
		if nci.Id == ci.Id {
			realCi = nci
		}
	}
	editW := storagejson.AppRef.A.NewWindow("编辑归类")
	vBox := container.NewVBox()
	vBox.Add(widget.NewLabel("名称："))
	nameEntry := widget.NewEntry()
	nameEntry.Text = realCi.Name
	vBox.Add(nameEntry)
	vBox.Add(widget.NewLabel("描述："))
	descriptionEntry := widget.NewEntry()
	descriptionEntry.Text = realCi.Description
	vBox.Add(descriptionEntry)
	vBox.Add(widget.NewLabel("Alias："))
	aliasEntry := widget.NewEntry()
	aliasEntry.Text = realCi.Alias
	vBox.Add(aliasEntry)
	// 创建水平按钮布局
	editCancelBtn := widget.NewButton("取消", func() {
		editW.Close()
	})
	editConfirmBtn := widget.NewButton("确定", func() {
		e := &common.EditForm{Name: nameEntry.Text, Alias: aliasEntry.Text, Description: descriptionEntry.Text}
		editConfirm(e, realCi, ciCard)
		editW.Close()
	})
	editConfirmBtn.Importance = widget.HighImportance
	hBox := container.NewHBox(widget.NewToolbarSpacer().ToolbarObject(), editCancelBtn, editConfirmBtn)
	vBox.Add(hBox)
	editW.SetContent(vBox)
	if !fyne.CurrentDevice().IsMobile() {
		editW.Resize(fyne.Size{Width: 300, Height: vBox.Size().Height})
	}
	editW.CenterOnScreen()
	editW.Show()
}

func editConfirm(e *common.EditForm, realCi storagejson.Category, ciCard *widget.Card) {
	if common.IsWhiteAndSpace(e.Name) {
		dialog.ShowInformation("提示", "归类文件夹名称不能是空的。", storagejson.AppRef.W)
		return
	}
	storagejson.EditCategory(e, realCi, ciCard)
}
