package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"pasecret/core/storagejson"
)

var theData storagejson.Data = storagejson.Data{}
var isEditOp bool

// showDataEditWin 点击了某密码项按钮，显示此具体编辑窗口
func showDataEditWin(performDataOrg *storagejson.Data, cidOrg string) {
	var editW fyne.Window
	if performDataOrg == nil {
		// 不是编辑，而是点击添加新增一个密码项
		isEditOp = false
		editW = storagejson.AppRef.A.NewWindow("添加账户密码项")
		theData.Name = ""
		theData.AccountName = ""
		theData.Password = ""
		theData.Site = ""
		theData.Remark = ""
		theData.CategoryId = cidOrg
		r, dId := common.GenAscRankId()
		if !r {
			dialog.ShowInformation("err", "showDataEditWin->common.GenAscRankId", storagejson.AppRef.W)
			return
		}
		theData.Id = dId
	} else {
		// 是编辑，填充当前密码项
		isEditOp = true
		editW = storagejson.AppRef.A.NewWindow("账户密码详细")
		theData.Name = (*performDataOrg).Name
		theData.AccountName = (*performDataOrg).AccountName
		theData.Password = (*performDataOrg).Password
		theData.Site = (*performDataOrg).Site
		theData.Remark = (*performDataOrg).Remark
		//the same: theData.categoryId = cidOrg
		theData.CategoryId = (*performDataOrg).CategoryId
		theData.Id = (*performDataOrg).Id

	}

	vBox := container.NewVBox()
	vBox.Add(widget.NewLabel("名称："))
	nameEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(nameEntry)
	nameEntry.Text = theData.Name
	vBox.Add(nameEntry)
	vBox.Add(widget.NewLabel("账号："))
	accountNameEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(accountNameEntry)
	accountNameEntry.Text = theData.AccountName
	vBox.Add(accountNameEntry)
	vBox.Add(widget.NewLabel("密码："))
	passwordEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(passwordEntry)
	passwordEntry.Text = theData.Password
	vBox.Add(passwordEntry)
	vBox.Add(widget.NewLabel("网址："))
	siteEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(siteEntry)
	siteEntry.Text = theData.Site
	vBox.Add(siteEntry)
	vBox.Add(widget.NewLabel("备注："))
	remarkEntry := widget.NewMultiLineEntry()
	common.EntryOnChangedEventHandler(remarkEntry)
	remarkEntry.Text = theData.Remark
	vBox.Add(remarkEntry)
	// 创建水平按钮布局
	copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		// 复制内容到系统剪切板
		copyStr := fmt.Sprintf("名称：\n%s\n账号：\n%s\n密码：\n%s\n网址：\n%s\n备注：\n%s",
			theData.Name, theData.AccountName, theData.Password, theData.Site, theData.Remark)
		editW.Clipboard().SetContent(copyStr)
		storagejson.AppRef.A.SendNotification(&fyne.Notification{
			Title:   "Pasecret",
			Content: "内容已复制到剪切板。",
		})
	})
	editCancelBtn := widget.NewButton("取消", func() {
		editW.Close()
	})
	editConfirmBtn := widget.NewButton("确定", func() {
		// 填充当前用户输入的值
		theData.Name = nameEntry.Text
		theData.AccountName = accountNameEntry.Text
		theData.Password = passwordEntry.Text
		theData.Site = siteEntry.Text
		theData.Remark = remarkEntry.Text
		editConfirmData(isEditOp, cidOrg, editW)
	})
	editConfirmBtn.Importance = widget.HighImportance
	var hBox *fyne.Container
	// 是编辑 则显示复制按钮
	if isEditOp {
		hBox = container.NewHBox(copyBtn, widget.NewToolbarSpacer().ToolbarObject(), editCancelBtn, editConfirmBtn)
	} else {
		hBox = container.NewHBox(widget.NewToolbarSpacer().ToolbarObject(), editCancelBtn, editConfirmBtn)
	}
	vBox.Add(hBox)
	editW.SetContent(vBox)
	if !fyne.CurrentDevice().IsMobile() {
		editW.Resize(fyne.Size{Width: 300, Height: vBox.Size().Height})
	}
	editW.CenterOnScreen()
	editW.Show()
}

func editConfirmData(isEditOp bool, cidOrg string, editW fyne.Window) {
	if common.IsWhiteAndSpace(theData.Name) {
		dialog.ShowInformation("提示", "名称不能是空的。", editW)
		return
	}
	storagejson.EditData(theData, isEditOp, cidOrg)
	editW.Close()
}
