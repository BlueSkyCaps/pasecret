package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"pasecret/core/storagejson"
	"unicode/utf8"
)

// ShowCategoryAddWin 首页点击了添加按钮，添加新归类夹
func ShowCategoryAddWin() {
	addW := storagejson.AppRef.A.NewWindow("编辑归类")
	vBox := container.NewVBox()
	vBox.Add(widget.NewLabel("名称："))
	nameEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(nameEntry)
	vBox.Add(nameEntry)
	vBox.Add(widget.NewLabel("描述："))
	descriptionEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(descriptionEntry)
	vBox.Add(descriptionEntry)
	vBox.Add(widget.NewLabel("Alias："))
	aliasEntry := widget.NewEntry()
	common.EntryOnChangedEventHandler(aliasEntry)
	vBox.Add(aliasEntry)
	// 创建水平按钮布局
	editCancelBtn := widget.NewButton("取消", func() {
		addW.Close()
	})
	editConfirmBtn := widget.NewButton("确定", func() {
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
		dialog.ShowInformation("提示", "归类文件夹名称不能是空的。", storagejson.AppRef.W)
		return
	}
	tips := ""
	if utf8.RuneCountInString(e.Name) > 12 {
		tips = "名称大于建议的长度:10字\n"
	}
	if utf8.RuneCountInString(e.Description) > 24 {
		tips = tips + "描述大于建议的长度:24字\n"
	}
	tips = tips + "\n是否保存？"
	dialog.ShowConfirm("提示", tips, func(b bool) {
		if b {
			addAsyncHandler(e)
			addW.Close()
		}
	}, addW)
}

func addAsyncHandler(e *common.EditForm) {
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
