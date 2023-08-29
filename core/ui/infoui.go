package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
	"pasecret/core/pi18n"
	"pasecret/core/storagedata"
)

// ShowCategoryInfoWin 点击了某归类文件夹Card的info按钮，显示此具体详情窗口
func ShowCategoryInfoWin(ci storagedata.Category) {
	var realCi storagedata.Category
	for _, nci := range storagedata.AppRef.LoadedItems.Category {
		if nci.Id == ci.Id {
			realCi = nci
		}
	}
	infoW := storagedata.AppRef.A.NewWindow("详情")
	vBox := container.NewVBox()
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("categoryInsetNameLabel", nil)))
	vBox.Add(canvas.NewText(realCi.Name, colornames.Lightgreen))
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("categoryInsetDescriptionLabel", nil)))
	vBox.Add(canvas.NewText(realCi.Description, colornames.Lightgreen))
	vBox.Add(widget.NewLabel("Alias："))
	vBox.Add(canvas.NewText(realCi.Alias, colornames.Lightgreen))
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("categoryInfoCanDelIfLabel", nil)))
	var re string
	if realCi.Removable {
		re = pi18n.LocalizedText("categoryInfoCanDelLabel", nil)
	} else {
		re = pi18n.LocalizedText("categoryInfoCanNotDelLabel", nil)
	}
	vBox.Add(canvas.NewText(re, colornames.Lightgreen))
	vBox.Add(widget.NewLabel(pi18n.LocalizedText("categoryInfoCanEditIfLabel", nil)))
	if realCi.Renameable {
		re = pi18n.LocalizedText("categoryInfoCanEditLabel", nil)
	} else {
		re = pi18n.LocalizedText("categoryInfoCanNotEditLabel", nil)
	}
	vBox.Add(canvas.NewText(re, colornames.Lightgreen))
	vBox.Add(widget.NewButton(pi18n.LocalizedText("categoryInsetOkButtonText", nil), func() {
		infoW.Close()
	}))
	infoW.SetContent(vBox)
	if !fyne.CurrentDevice().IsMobile() {
		infoW.Resize(fyne.Size{Width: 300, Height: vBox.Size().Height})
	}
	infoW.CenterOnScreen()
	infoW.Show()
}
