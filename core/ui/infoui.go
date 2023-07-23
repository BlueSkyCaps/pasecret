package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
	"pasecret/core/storagejsondata"
)

// ShowCategoryInfoWin 点击了某归类文件夹Card的info按钮，显示此具体详情窗口
func ShowCategoryInfoWin(ci storagejsondata.Category) {
	infoW := storagejsondata.AppRef.A.NewWindow("详情")
	vBox := container.NewVBox()
	vBox.Add(widget.NewLabel("名称："))
	vBox.Add(canvas.NewText(ci.Name, colornames.Darkblue))
	vBox.Add(widget.NewLabel("描述："))
	vBox.Add(canvas.NewText(ci.Description, colornames.Darkblue))
	vBox.Add(widget.NewLabel("Alias："))
	vBox.Add(canvas.NewText(ci.Alias, colornames.Darkblue))
	vBox.Add(widget.NewLabel("存储的密码项："))
	vBox.Add(canvas.NewText("好多个", colornames.Darkblue))
	vBox.Add(widget.NewLabel("可被删除："))
	var re string
	if ci.Removable {
		re = "可以删除"
	} else {
		re = "内置归类夹，无法删除。"
	}
	vBox.Add(canvas.NewText(re, colornames.Darkblue))
	vBox.Add(widget.NewLabel("可被重命名："))
	if ci.Renameable {
		re = "可以重命名"
	} else {
		re = "内置归类夹，无法重命名。"
	}
	vBox.Add(canvas.NewText(re, colornames.Darkblue))
	vBox.Add(widget.NewButton("确定", func() {
		infoW.Close()
	}))
	infoW.SetContent(vBox)
	if !fyne.CurrentDevice().IsMobile() {
		infoW.Resize(fyne.Size{Width: 300, Height: vBox.Size().Height})
	}
	infoW.CenterOnScreen()
	infoW.Show()
}
