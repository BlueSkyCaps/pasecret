package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"pasecret/core/storagejsondata"
)

var (
	name        string
	accountName string
	password    string
	site        string
	remark      string
	categoryId  string
)

// showDataEditWin 点击了某密码项按钮，显示此具体编辑窗口
func showDataEditWin(performData *storagejson.Data) {

	if performData == nil {
		name = ""
		accountName = ""
		password = ""
		site = ""
		remark = ""
		categoryId = ""
	} else {
		name = (*performData).Name
		accountName = (*performData).AccountName
		password = (*performData).Password
		site = (*performData).Site
		remark = (*performData).Remark
		categoryId = (*performData).CategoryId
	}

	editW := storagejson.AppRef.A.NewWindow("账户密码详细")
	vBox := container.NewVBox()
	vBox.Add(widget.NewLabel("名称："))
	nameEntry := widget.NewEntry()
	nameEntry.Text = name
	vBox.Add(nameEntry)
	vBox.Add(widget.NewLabel("账号名字："))
	accountNameEntry := widget.NewEntry()
	accountNameEntry.Text = accountName
	vBox.Add(accountNameEntry)
	vBox.Add(widget.NewLabel("密码："))
	passwordEntry := widget.NewEntry()
	passwordEntry.Text = password
	vBox.Add(passwordEntry)
	vBox.Add(widget.NewLabel("网址："))
	siteEntry := widget.NewEntry()
	siteEntry.Text = site
	vBox.Add(siteEntry)
	vBox.Add(widget.NewLabel("备注："))
	remarkEntry := widget.NewEntry()
	remarkEntry.Text = remark
	vBox.Add(remarkEntry)
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

func editConfirmData(e *common.EditForm, realCi storagejson.Category, ciCard *widget.Card) {
	if common.IsWhiteAndSpace(e.Name) {
		dialog.ShowInformation("提示", "归类文件夹名称不能是空的。", storagejson.AppRef.W)
		return
	}
	storagejson.EditCategory(e, realCi, ciCard)
}
