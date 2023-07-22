// Package ui 定义窗口主要UI内容元素
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
	"sync"
	"time"
)

var w fyne.Window
var a fyne.App
var searchInput *widget.Entry
var searchBtn *widget.Button

// Run 开始定义UI元素并显示窗口
func Run(tw fyne.Window, ta fyne.App) {
	w = tw
	a = ta
	// 创建工具条容器，包含了添加按钮、搜索按钮
	toolbarBox := createHBox()
	// 创建网格容器父布局
	gridParent := container.NewMax()
	// 创建网格容器
	var grid *fyne.Container
	if !fyne.CurrentDevice().IsMobile() {
		grid = container.NewGridWrap(fyne.Size{Width: 200, Height: 150})
	} else {
		grid = container.NewGridWithColumns(2)
	}
	// 加载密码归类文件夹项到网格容器
	grid.Objects = loadItems()
	// 加载网格容器到它的父布局
	gridParent.Add(grid)
	// 将工具条和文件夹网格列表放进border，形成垂直布局效果
	homeTab := container.NewBorder(toolbarBox, nil, nil, nil, container.NewVScroll(gridParent))
	// 添加tabs选项卡
	appTabs := firstAddTabs(homeTab)
	// 设置窗体最终布局内容
	w.SetContent(appTabs)
	w.ShowAndRun()
}

func loadItems() []fyne.CanvasObject {
	var items []fyne.CanvasObject
	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(i_ int, wg_ *sync.WaitGroup) {
			defer wg_.Done()
			currentCart := createCurrentCart(strconv.Itoa(i_))
			items = append(items, currentCart)
		}(i, &wg)
		time.Sleep(time.Millisecond * 10)
	}
	wg.Wait()

	return items
}

func createCurrentCart(i string) *widget.Card {
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.MoreHorizontalIcon(), func() {
			log.Println("New document")
		}),
		widget.NewToolbarAction(theme.InfoIcon(), func() {}),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {}),
	)
	card := widget.NewCard("cart"+i, "文件夹文件夹文件夹"+i, toolbar)
	return card
}

func firstAddTabs(home *fyne.Container) *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.HomeIcon(), home),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), widget.NewButton("Open new", func() {
			w3 := a.NewWindow("Third")
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
		dialog.ShowInformation("查找", searchInput.Text, w)
	})
	vBoxLayout := container.NewHBox(addMenuToolbar(), searchInput, searchBtn)
	return vBoxLayout
}

func createGrid() fyne.Layout {
	gridLayout := layout.NewGridLayout(3)
	return gridLayout
}
