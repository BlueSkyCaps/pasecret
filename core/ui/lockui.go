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

func LockUI() fyne.Window {
	currentNumberClickCount = 0
	currentValidNumbs = ""
	storagejson.AppRef.LockWin = storagejson.AppRef.A.NewWindow("Pasecret")
	if fyne.CurrentDevice().IsMobile() {
		storagejson.AppRef.LockWin.SetOnClosed(func() {
			// 关闭解锁窗口回调，立马重新生成显示解锁窗口，因为上一次的解锁窗口资源已被释放，必须重新调用生成
			LockUI().Show()
		})
	}

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
	storagejson.AppRef.LockWin.SetContent(center)
	return storagejson.AppRef.LockWin
}

// 点击了任意数字按钮
func enterNumberBtnHandler(i string) {
	// 更新label圆圈表示当前已键入
	currentNumberClickCount = currentNumberClickCount + 1
	displayNumberLabels[currentNumberClickCount-1].SetText(displayInChar)
	// 追加当前点击的数字
	currentValidNumbs = currentValidNumbs + i
	// 若点击了第四次，则进行验证
	if len(currentValidNumbs) == 4 {
		lockpw := preferences.GetPreferenceByLockPwd()
		if lockpw == currentValidNumbs {
			if !fyne.CurrentDevice().IsMobile() {
				storagejson.AppRef.W.Show()
				storagejson.AppRef.LockWin.Close()
			} else {
				// 安卓端必须隐藏窗口来代替Close关闭窗口，因为Close会递归引发SetOnClosed事件
				storagejson.AppRef.LockWin.Hide()
				// 充值本次窗口回调，再关闭窗口，不会引起原事件
				storagejson.AppRef.LockWin.SetOnClosed(func() {
				})
				storagejson.AppRef.LockWin.Close()
			}
		}
		// 验证失败，重新改变状态为初始化
		currentValidNumbs = ""
		currentNumberClickCount = 0
		for i := 0; i < len(displayNumberLabels); i++ {
			displayNumberLabels[i].SetText(displayUnChar)
		}
	}
}
