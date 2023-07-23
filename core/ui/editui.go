package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/storagejsondata"
)

// ShowCategoryEditWin 点击了某归类文件夹Card的Edit按钮，显示此具体编辑窗口
func ShowCategoryEditWin(ci storagejsondata.Category) {
	editW := storagejsondata.AppRef.A.NewWindow("编辑归类")
	vBox := container.NewVBox()
	vBox.Add(widget.NewLabel("名称："))
	nameEntry := widget.NewEntry()
	nameEntry.Text = ci.Name
	vBox.Add(nameEntry)
	vBox.Add(widget.NewLabel("描述："))
	descriptionEntry := widget.NewEntry()
	descriptionEntry.Text = ci.Description
	vBox.Add(descriptionEntry)
	vBox.Add(widget.NewLabel("Alias："))
	aliasEntry := widget.NewEntry()
	aliasEntry.Text = ci.Alias
	vBox.Add(aliasEntry)
	// 创建水平按钮布局
	editCancelBtn := widget.NewButton("取消", func() {
		editW.Close()
	})
	editConfirmBtn := widget.NewButton("确定", func() {
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
