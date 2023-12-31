package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"pasecret/core/pi18n"
	"pasecret/core/preferences"
	"pasecret/core/storagedata"
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
	storagedata.AppRef.LockWin = storagedata.AppRef.A.NewWindow("Pasecret")
	if fyne.CurrentDevice().IsMobile() {
		storagedata.AppRef.LockWin.SetOnClosed(func() {
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
	vBox.Add(widget.NewLabelWithStyle(pi18n.LocalizedText("lockLabelText", nil),
		fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	vBox.Add(displayNumberHBox)
	vBox.Add(enterNumberGrid)

	center := container.NewCenter()
	center.Add(vBox)
	storagedata.AppRef.LockWin.SetContent(center)
	return storagedata.AppRef.LockWin
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
		lockPwdPure, err := common.DecryptAES([]byte(common.AppProductKeyAES), preferences.GetPreferenceByLockPwd())
		if err != nil {
			dialog.ShowError(err, storagedata.AppRef.LockWin)
			return
		}
		if currentValidNumbs == lockPwdPure {
			if !fyne.CurrentDevice().IsMobile() {
				storagedata.AppRef.W.Show()
				storagedata.AppRef.LockWin.Close()
			} else {
				// 安卓端必须隐藏窗口来代替Close关闭窗口，因为Close会递归引发SetOnClosed事件
				storagedata.AppRef.LockWin.Hide()
				// 重置本次窗口回调，再关闭窗口，不会引起原事件
				storagedata.AppRef.LockWin.SetOnClosed(func() {
				})
				storagedata.AppRef.LockWin.Close()
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
