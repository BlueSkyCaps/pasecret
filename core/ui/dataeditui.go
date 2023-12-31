package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"pasecret/core/pi18n"
	"pasecret/core/storagedata"
)

var theData storagedata.Data = storagedata.Data{}
var isEditOp bool

// showDataEditWin 点击了某密码项按钮，显示此具体编辑窗口
func showDataEditWin(performDataOrg *storagedata.Data, cidOrg string) {
	var editW fyne.Window
	if performDataOrg == nil {
		// 不是编辑，而是点击添加新增一个密码项
		isEditOp = false
		editW = storagedata.AppRef.A.NewWindow(pi18n.LocalizedText("dataInsertWindowTitle", nil))
		theData.Name = ""
		theData.AccountName = ""
		theData.Password = ""
		theData.Site = ""
		theData.Remark = ""
		theData.CategoryId = cidOrg
		r, dId := common.GenAscRankId()
		if !r {
			dialog.ShowInformation("err", "showDataEditWin->common.GenAscRankId", storagedata.AppRef.W)
			return
		}
		theData.Id = dId
	} else {
		// 是编辑，填充当前密码项
		isEditOp = true
		editW = storagedata.AppRef.A.NewWindow(pi18n.LocalizedText("dataEditWindowTitle", nil))
		theData.Name = (*performDataOrg).Name
		theData.AccountName = (*performDataOrg).AccountName
		theData.Password = (*performDataOrg).Password
		theData.Site = (*performDataOrg).Site
		theData.Remark = (*performDataOrg).Remark
		//the same: theData.categoryId = cidOrg
		theData.CategoryId = (*performDataOrg).CategoryId
		theData.Id = (*performDataOrg).Id

	}
	accountLabelText := dealEditDataAccountLabelText(theData.CategoryId)
	vBox := container.NewVBox()
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("dataEditNameLabel", nil)))
	nameEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(nameEntry)
	nameEntry.Text = theData.Name
	vBox.Add(nameEntry)

	vBox.Add(widget.NewLabel(accountLabelText))
	accountNameEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(accountNameEntry)
	accountNameEntry.Text = theData.AccountName
	vBox.Add(accountNameEntry)
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("dataEditPwdLabel", nil)))
	passwordEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(passwordEntry)
	passwordEntry.Text = theData.Password
	vBox.Add(passwordEntry)
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("dataEditSiteLabel", nil)))
	siteEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(siteEntry)
	siteEntry.Text = theData.Site
	vBox.Add(siteEntry)
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("dataEditRemarkLabel", nil)))
	remarkEntry := widget.NewMultiLineEntry()
	common.EntryOnChangedEventHandler(remarkEntry)
	remarkEntry.Text = theData.Remark
	vBox.Add(remarkEntry)
	// 创建水平按钮布局
	copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		// 复制内容到系统剪切板
		copyStr := pi18n.LocalizedText("copyStr", map[string]interface{}{
			"name":        theData.Name,
			"accountName": theData.AccountName,
			"password":    theData.Password,
			"site":        theData.Site,
			"remark":      theData.Remark,
		})
		editW.Clipboard().SetContent(copyStr)
		storagedata.AppRef.A.SendNotification(&fyne.Notification{
			Title:   "Pasecret",
			Content: pi18n.LocalizedText("copyClipboardDone", nil),
		})
	})
	editCancelBtn := widget.NewButton(pi18n.LocalizedText("dataEditCancelButtonText", nil), func() {
		editW.Close()
	})
	editConfirmBtn := widget.NewButton(pi18n.LocalizedText("dataEditOkButtonText", nil), func() {
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

// 根据是哪个归类夹id定义账号标签文本，如内置的银行卡，可以显示"卡号"而不是默认"账号"
func dealEditDataAccountLabelText(cid string) string {
	switch cid[0] {
	case '0':
		return pi18n.LocalizedText("dataEditDynamicByIdNameLabel", nil)
	case '1':
		return pi18n.LocalizedText("dataEditDynamicByBankNameLabel", nil)
	default:
		return pi18n.LocalizedText("dataEditDynamicByUsualAccountNameLabel", nil)
	}
}

func editConfirmData(isEditOp bool, cidOrg string, editW fyne.Window) {
	if common.IsWhiteAndSpace(theData.Name) {
		dialog.ShowInformation(pi18n.LocalizedText("dialogShowInformationTitle", nil),
			pi18n.LocalizedText("dataEditNameBlank", nil), editW)
		return
	}
	storagedata.EditData(theData, isEditOp, cidOrg)
	editW.Close()
}
