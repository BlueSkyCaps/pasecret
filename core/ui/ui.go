// Package ui 定义窗口主要UI内容元素
package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"pasecret/core/storagejsondata"
	"sync"
	"time"
)

var appRef storagejsondata.AppStructRef
var searchInput *widget.Entry
var searchBtn *widget.Button

// Run 开始定义UI元素并显示窗口
func Run(_appRef storagejsondata.AppStructRef) {
	appRef = _appRef
	// 创建工具条容器，包含了添加按钮、搜索按钮
	toolbarBox := createHBox()
	// 创建网格容器父布局
	gridParent := container.NewMax()
	// 创建网格容器
	grid := createGrid()
	// 加载密码归类文件夹项到网格容器
	grid.Objects = loadItems()
	// 加载网格容器到它的父布局
	gridParent.Add(grid)
	// 将工具条和文件夹网格列表放进border，形成垂直布局效果
	homeTab := container.NewBorder(toolbarBox, nil, nil, nil, container.NewVScroll(gridParent))
	// 添加tabs选项卡
	appTabs := firstAddTabs(homeTab)
	// 设置窗体最终布局内容
	appRef.W.SetContent(appTabs)
	appRef.W.ShowAndRun()
}

func loadItems() []fyne.CanvasObject {
	var cardItemsF []fyne.CanvasObject
	var wg sync.WaitGroup
	categoryLen := len(appRef.LoadedItems.Category)
	wg.Add(categoryLen)

	for i := 0; i < categoryLen; i++ {
		go func(i_ int, wg_ *sync.WaitGroup) {
			defer wg_.Done()
			currentCart := createCurrentCart(appRef.LoadedItems.Category[i])
			cardItemsF = append(cardItemsF, currentCart)
		}(i, &wg)
		time.Sleep(time.Millisecond * 10)
	}
	wg.Wait()
	return cardItemsF
}

func createCurrentCart(i storagejsondata.Category) *widget.Card {
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.MoreHorizontalIcon(), func() {

			dialog.ShowInformation("w", i.Name, appRef.W)
		}),
		widget.NewToolbarAction(theme.InfoIcon(), func() {

		}),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			if !i.Removable {
				dialog.ShowInformation("提示", "该归类不可被删除。", appRef.W)
				return
			}
			dialog.ShowConfirm("提示", fmt.Sprintf("确定删除%s？该归类下保存的所有密码一并会被删除。", i.Name),
				func(b bool) {
					if b {
						dialog.ShowInformation("提示", "删了。", appRef.W)
						return
					}
					dialog.ShowInformation("提示", "点了不删。", appRef.W)
				}, appRef.W)
		}),
	)
	card := widget.NewCard(i.Name, i.Description, toolbar)
	return card
}

func firstAddTabs(home *fyne.Container) *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.HomeIcon(), home),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), widget.NewButton("Open new", func() {
			w3 := appRef.A.NewWindow("Third")
			w3.SetContent(widget.NewButton("确定", func() {
				w3.Close()
			}))

			w3.Show()
		})),
	)
	tabs.SetTabLocation(container.TabLocationBottom)
	return tabs
}

func addMenuToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			log.Println("New document")
		}),
	)
	return toolbar
}

func createHBox() *fyne.Container {
	searchInput = widget.NewEntry()
	searchBtn = widget.NewButtonWithIcon("查找", theme.SearchIcon(), func() {
		dialog.ShowInformation("查找", searchInput.Text, appRef.W)
	})
	vBoxLayout := container.NewHBox(addMenuToolbar(), searchInput, searchBtn)
	return vBoxLayout
}

func createGrid() *fyne.Container {
	var grid *fyne.Container
	if !fyne.CurrentDevice().IsMobile() {
		grid = container.NewGridWrap(fyne.Size{Width: 200, Height: 150})
	} else {
		grid = container.NewGridWithColumns(2)
	}
	return grid
}
