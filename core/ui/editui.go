package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"pasecret/core/pi18n"
	"pasecret/core/storagedata"
	"unicode/utf8"
)

// ShowCategoryEditWin 点击了某归类文件夹Card的Edit按钮，显示此具体编辑窗口
func ShowCategoryEditWin(ci storagedata.Category, ciCard *widget.Card) {
	var realCi storagedata.Category
	// 从AppRef中根据id找到文件夹，因为回调函数传的Category参数是最初原始的数据，而AppRef中LoadedItems是实时更新的数据
	for _, nci := range storagedata.AppRef.LoadedItems.Category {
		if nci.Id == ci.Id {
			realCi = nci
		}
	}
	editW := storagedata.AppRef.A.NewWindow(pi18n.LocalizedText("categoryInsetWindowTitle", nil))
	vBox := container.NewVBox()
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("categoryInsetNameLabel", nil)))
	nameEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(nameEntry)
	nameEntry.Text = realCi.Name
	vBox.Add(nameEntry)
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("categoryInsetDescriptionLabel", nil)))
	descriptionEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(descriptionEntry)
	descriptionEntry.Text = realCi.Description
	vBox.Add(descriptionEntry)
	vBox.Add(widget.NewLabel("Alias："))
	aliasEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(aliasEntry)
	aliasEntry.Text = realCi.Alias
	vBox.Add(aliasEntry)
	// 创建水平按钮布局
	editCancelBtn := widget.NewButton(pi18n.LocalizedText("categoryInsetCancelButtonText", nil), func() {
		editW.Close()
	})
	editConfirmBtn := widget.NewButton(pi18n.LocalizedText("categoryInsetOkButtonText", nil), func() {
		e := &common.EditForm{Name: nameEntry.Text, Alias: aliasEntry.Text, Description: descriptionEntry.Text}
		editConfirm(e, realCi, ciCard, editW)
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

func editConfirm(e *common.EditForm, realCi storagedata.Category, ciCard *widget.Card, editW fyne.Window) {
	if common.IsWhiteAndSpace(e.Name) {
		dialog.ShowInformation(pi18n.LocalizedText("dialogShowInformationTitle", nil),
			pi18n.LocalizedText("categoryInsetNameBlank", nil), editW)
		return
	}
	tips := ""
	if utf8.RuneCountInString(e.Name) > 12 {
		tips = pi18n.LocalizedText("categoryInsetNameLength", nil)
	}
	if utf8.RuneCountInString(e.Description) > 24 {
		tips = tips + pi18n.LocalizedText("categoryInsetDescriptionLength", nil)
	}
	tips = tips + pi18n.LocalizedText("categoryInsetNameConfirm", nil)

	dialog.ShowConfirm(pi18n.LocalizedText("dialogShowInformationTitle", nil), tips, func(b bool) {
		if b {
			editAsyncHandler(e, realCi, ciCard)
			editW.Close()
		}
	}, editW)

}

func editAsyncHandler(e *common.EditForm, realCi storagedata.Category, ciCard *widget.Card) {
	storagedata.EditCategory(e, realCi, ciCard)
}
