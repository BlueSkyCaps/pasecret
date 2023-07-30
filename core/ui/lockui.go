package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/preferences"
	"pasecret/core/storagejson"
	"strconv"
)

const displayUnChar = "    ○    "
const displayInChar = "    ●    "

var displayNumberLabels []*widget.Label = []*widget.Label{
	widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
}
var currentNumberClickCount uint8
var currentValidNumbs string
var lockWin fyne.Window

func LockUI() fyne.Window {
	currentNumberClickCount = 0
	currentValidNumbs = ""
	lockWin = storagejson.AppRef.A.NewWindow("Pasecret")
	displayNumberHBox := container.NewHBox()
	for i := 0; i < len(displayNumberLabels); i++ {
		displayNumberLabels[i].SetText(displayUnChar)
		displayNumberHBox.Add(displayNumberLabels[i])
	}
	enterNumberGrid := container.NewGridWithColumns(3)
	for i := 1; i <= 10; i++ {
		if i == 10 {
			enterNumberGrid.Add(widget.NewToolbarSpacer().ToolbarObject())
			button := widget.NewButton("0", func() {
				enterNumberBtnHandler("0")
			})
			enterNumberGrid.Add(button)
			enterNumberGrid.Add(widget.NewToolbarSpacer().ToolbarObject())
			break
		}
		i := i
		button := widget.NewButton(strconv.Itoa(i), func() {
			enterNumberBtnHandler(strconv.Itoa(i))
		})
		enterNumberGrid.Add(button)
	}
	vBox := container.NewVBox()
	vBox.Add(widget.NewLabelWithStyle("解锁", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	vBox.Add(displayNumberHBox)
	vBox.Add(enterNumberGrid)

	center := container.NewCenter()
	center.Add(vBox)
	lockWin.SetContent(center)
	return lockWin
}

// 点击了任意数字按钮
func enterNumberBtnHandler(i string) {
	// 更新label圆圈表示当前已键入
	currentNumberClickCount = currentNumberClickCount + 1
	displayNumberLabels[currentNumberClickCount-1].SetText(displayInChar)
	// 追加当前点击的数字
	currentValidNumbs = currentValidNumbs + i
	// 若点击了第四次，则进行验证
	println(currentValidNumbs)
	if len(currentValidNumbs) == 4 {
		lockpw := preferences.GetPreferenceByLockPwd()
		if lockpw == currentValidNumbs {
			// 打开主窗口，关闭解锁窗口
			storagejson.AppRef.W.Show()
			lockWin.Close()
		}
		// 验证失败，重新改变状态为初始化
		currentValidNumbs = ""
		currentNumberClickCount = 0
		for i := 0; i < len(displayNumberLabels); i++ {
			displayNumberLabels[i].SetText(displayUnChar)
		}
	}
}
