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

// ShowCategoryAddWin 首页点击了添加按钮，添加新归类夹
func ShowCategoryAddWin() {
	addW := storagedata.AppRef.A.NewWindow(pi18n.LocalizedText("categoryInsetWindowTitle", nil))
	vBox := container.NewVBox()
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("categoryInsetNameLabel", nil)))
	nameEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(nameEntry)
	vBox.Add(nameEntry)
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("categoryInsetDescriptionLabel", nil)))
	descriptionEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(descriptionEntry)
	vBox.Add(descriptionEntry)
	vBox.Add(widget.NewLabel("Alias："))
	aliasEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(aliasEntry)
	vBox.Add(aliasEntry)
	// 创建水平按钮布局
	editCancelBtn := widget.NewButton(pi18n.LocalizedText("categoryInsetCancelButtonText", nil),
		func() {
			addW.Close()
		})
	editConfirmBtn := widget.NewButton(pi18n.LocalizedText("categoryInsetOkButtonText", nil),
		func() {
			e := &common.EditForm{Name: nameEntry.Text, Alias: aliasEntry.Text, Description: descriptionEntry.Text}
			AddConfirm(e, addW)
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

func AddConfirm(e *common.EditForm, addW fyne.Window) {
	if common.IsWhiteAndSpace(e.Name) {
		dialog.ShowInformation(pi18n.LocalizedText("dialogShowInformationTitle", nil),
			pi18n.LocalizedText("categoryInsetNameBlank", nil), addW)
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
			addAsyncHandler(e)
			addW.Close()
		}
	}, addW)
}

func addAsyncHandler(e *common.EditForm) {
	// 先成功新增到本地存储库
	addCi := storagedata.AddCategory(e)
	if addCi == nil {
		return
	}
	// 再根据Category对象创建引用Cart小部件
	cart := CreateCurrentCart(*addCi)
	// 最后更新当前Cart小部件布局
	storagedata.AppRef.RepaintCartsByAdd(*addCi, cart)
}
