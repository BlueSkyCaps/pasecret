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
	"sync"
	"time"
)

// Run 开始定义UI元素并显示窗口
func Run() {
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
	storagejson.AppRef.W.ShowAndRun()
}

// AppTabsRefreshHandler AppTabs的第二个tab页（设置页）在安卓端后台运行一段时间重新打开（基于安卓后台贮存机制不会关闭应用），
// 设置页会显示空白。此处开启个go重新渲染
func AppTabsRefreshHandler(tabs *container.AppTabs) {
	if !fyne.CurrentDevice().IsMobile() {
		return
	}
	go func() {
		for true {
			time.Sleep(time.Second * 2)
			if tabs.SelectedIndex() == 1 {
				tabs.Items[1].Content.Refresh()
			}
		}
	}()
}

func loadItems() []fyne.CanvasObject {
	var cardItemsF []fyne.CanvasObject
	var wg sync.WaitGroup
	categoryLen := len(storagejson.AppRef.LoadedItems.Category)
	wg.Add(categoryLen)

	for i := 0; i < categoryLen; i++ {
		go func(i_ int, wg_ *sync.WaitGroup) {
			defer wg_.Done()
			currentCart := CreateCurrentCart(storagejson.AppRef.LoadedItems.Category[i])
			cardItemsF = append(cardItemsF, currentCart)
		}(i, &wg)
		time.Sleep(time.Millisecond * 10)
	}
	wg.Wait()
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
	storagejson.AppRef.SearchInput = widget.NewEntry()
	common.EntryOnChangedEventHandler(storagejson.AppRef.SearchInput)
	storagejson.AppRef.SearchBtn = widget.NewButtonWithIcon("查找", theme.SearchIcon(), func() {
		if common.IsWhiteAndSpace(storagejson.AppRef.SearchInput.Text) {
			return
		}
		ShowSearchResultWin()
	})
	vBoxLayout := container.NewHBox(addCategoryMenuToolbar(), storagejson.AppRef.SearchInput, storagejson.AppRef.SearchBtn)
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
