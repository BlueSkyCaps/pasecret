// Package ui 定义窗口主要UI内容元素
package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"pasecret/core/common"
	storagejson "pasecret/core/storagejson"
	"time"
)

// Run 开始定义UI元素并显示窗口
func Run(lock bool) {
	// 创建工具条容器，包含了添加按钮、搜索按钮
	toolbarBox := createHBox()
	// 创建网格容器父布局
	gridParent := container.NewMax()
	// 创建网格容器
	grid := createGrid()
	// 加载密码归类文件夹项到网格容器
	grid.Objects = loadItems()
	// 将具有密码归类文件夹项的网格容器指向AppRef
	storagejson.AppRef.CardsGrid = grid
	// 加载网格容器到它的父布局
	gridParent.Add(grid)
	// 将工具条和文件夹网格列表放进border，形成垂直布局效果
	homeTabContent := container.NewBorder(toolbarBox, nil, nil, nil, container.NewVScroll(gridParent))
	settingTabContent := container.NewBorder(nil, nil, nil, nil, createSettingTabContent())
	// 添加tabs选项卡
	appTabs := firstAddTabs(homeTabContent, settingTabContent)
	AppTabsRefreshHandler(appTabs)
	// 设置窗体最终布局内容
	storagejson.AppRef.W.SetContent(appTabs)
	if lock {
		if !fyne.CurrentDevice().IsMobile() {
			storagejson.AppRef.LockWin = LockUI()
			storagejson.AppRef.LockWin.Show()
			storagejson.AppRef.A.Run()
			return
		}
		/*注意 安卓端如果先显示解锁窗口，则后续解锁成功Show显示主窗口W会造成显示面积只有大约1/4。
		必须第一时间W.ShowAndRun()才能正常显示整个屏幕。因此采用开启一个goroutine，等待一会立马
		置顶显示解锁窗口，因为执行ShowAndRun()会阻塞后面代码。但点击解锁窗口关闭按钮，回调事件立马
		再次显示解锁窗口。只有在输入正确密码后才会主动隐藏解锁窗口*/
		go func() {
			time.Sleep(time.Second * 1)
			storagejson.AppRef.LockWin = LockUI()
			storagejson.AppRef.LockWin.Show()
		}()
		storagejson.AppRef.W.ShowAndRun()
	} else {
		storagejson.AppRef.W.Show()
		storagejson.AppRef.A.Run()
	}

}

// AppTabsRefreshHandler AppTabs的第二个tab页（设置页）运行一段时间重新显示（安卓后台贮存机制不会关闭应用），
// 设置页会显示空白。此处开启个go重新渲染
func AppTabsRefreshHandler(tabs *container.AppTabs) {
	go func() {
		for true {
			time.Sleep(time.Second * 1)
			if tabs.SelectedIndex() == 1 {
				tabs.Items[1].Content.Refresh()
			}
		}
	}()
}

func LockUILifecycleHandler() {
	if !fyne.CurrentDevice().IsMobile() {
		return
	}
	storagejson.AppRef.A.Lifecycle().SetOnExitedForeground(func() {
		println("lost")
	})
	storagejson.AppRef.A.Lifecycle().SetOnEnteredForeground(func() {
		println("get")
	})
}

func loadItems() []fyne.CanvasObject {
	var cardItemsF []fyne.CanvasObject
	categoryLen := len(storagejson.AppRef.LoadedItems.Category)

	for i := 0; i < categoryLen; i++ {
		currentCart := CreateCurrentCart(storagejson.AppRef.LoadedItems.Category[i])
		cardItemsF = append(cardItemsF, currentCart)
	}
	return cardItemsF
}

func CreateCurrentCart(ci storagejson.Category) *widget.Card {
	card := widget.NewCard(ci.Name, ci.Description, canvas.NewText("", color.RGBA{}))
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ListIcon(), func() {
			ShowDataList(ci)
		}),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			ShowCategoryInfoWin(ci)
		}),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			if !ci.Renameable {
				dialog.ShowInformation("提示", "该归类不可被编辑。", storagejson.AppRef.W)
				return
			}
			ShowCategoryEditWin(ci, card)
		}),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			if !ci.Removable {
				dialog.ShowInformation("提示", "该归类不可被删除。", storagejson.AppRef.W)
				return
			}
			dialog.ShowConfirm("提示", fmt.Sprintf("确定删除\n‘%s’？\n该归类保存的所有密码一并会被删除。", ci.Name),
				func(b bool) {
					if b {
						go func() {
							storagejson.DeleteCategory(ci, card)
						}()
						return
					}
				}, storagejson.AppRef.W)
		}),
	)
	card.SetContent(toolbar)
	return card
}

func firstAddTabs(home *fyne.Container, setting *fyne.Container) *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.HomeIcon(), home),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), setting),
	)
	tabs.SetTabLocation(container.TabLocationBottom)
	return tabs
}

func addCategoryMenuToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			ShowCategoryAddWin()
		}),
	)
	return toolbar
}

func createHBox() *fyne.Container {
	searchInputEntry := widget.NewEntry()
	// 安卓端 专为searchInputEntry控制的
	common.SearchEntryOnChangedEventHandler(searchInputEntry)
	searchBtn := widget.NewButtonWithIcon("查找", theme.SearchIcon(), func() {
		if common.IsWhiteAndSpace(searchInputEntry.Text) {
			return
		}
		ShowSearchResultWin(searchInputEntry.Text)
		/*安卓端 searchInputEntry的Text更改，common.SearchTmp也要同步更改，因为common.SearchTmp存储的还是是之前的值
		安卓端common.SearchEntryOnChangedEventHandler控制退格键避免退两次字符。
		windows忽略此问题*/
		// 在安卓端，必須Text賦值且Refresh才會更新文本框值,SetText無效
		searchInputEntry.SetText("")
		searchInputEntry.Text = ""
		common.SearchTmp = ""
		searchInputEntry.Refresh()
	})
	vBoxLayout := container.NewHBox(addCategoryMenuToolbar(), searchInputEntry, searchBtn)
	return vBoxLayout
}

func createGrid() *fyne.Container {
	var grid *fyne.Container
	if !fyne.CurrentDevice().IsMobile() {
		// NewGridWithColumns会自适应内容响应，如Grid由card组成，card标题不会换行，NewGridWithColumns可以自适应其宽度
		// NewGridWrap不会适应card的标题宽度，但可以缩放窗体自动填充能承载的最大列数
		grid = container.NewGridWithColumns(3)
		//grid = container.NewGridWrap(fyne.Size{Width: 200, Height: 150})
	} else {
		grid = container.NewGridWithColumns(1)
	}
	return grid
}
