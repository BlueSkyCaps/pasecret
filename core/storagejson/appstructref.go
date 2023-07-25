package storagejson

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"time"
)

// AppStructRef App操作句柄
type AppStructRef struct {
	W                   fyne.Window
	A                   fyne.App
	LoadedItems         LoadedItems
	SearchInput         *widget.Entry
	SearchBtn           *widget.Button
	CardsGrid           *fyne.Container
	DataList            *widget.List
	ShowDataEditWinFunc func(performData *Data, cid string)
}

func (appRef AppStructRef) RepaintCartsByRemove(delCard *widget.Card) {
	appRef.CardsGrid.Remove(delCard)
	appRef.CardsGrid.Refresh()
}

// RepaintCartsByEdit 成功保存本地存储库后再刷新Cart文件夹小部件
func (appRef AppStructRef) RepaintCartsByEdit(e *common.EditForm, editCard *widget.Card) {
	editCard.SetTitle(e.Name)
	editCard.SetSubTitle(e.Description)
	dialog.ShowInformation(editCard.Title, editCard.Subtitle, appRef.W)
	appRef.CardsGrid.Refresh()
	if fyne.CurrentDevice().IsMobile() {
		// 安卓端必须添加sleep阻塞一段时间才会重绘Cart文本，但是同理的添加删除Cart却能正常刷新显示
		time.Sleep(time.Millisecond * 500)
	}
}

// RepaintCartsByAdd 成功保存本地存储库后再刷新Cart文件夹小部件
func (appRef AppStructRef) RepaintCartsByAdd(addCi Category, addCard *widget.Card) {
	addCard.Title = addCi.Name
	addCard.Subtitle = addCi.Description
	appRef.CardsGrid.Objects = append(appRef.CardsGrid.Objects, addCard)
	appRef.CardsGrid.Refresh()
}

// RepaintDataListByEdit 成功保存本地存储库后再刷新List密码项列表，避免本地保存失败却事先更新列表
func (appRef AppStructRef) RepaintDataListByEdit(cidOrg string) {
	// 获取此文件夹的密码项 最新数据
	relatedData := GetRelatedDataByCid(cidOrg)
	appRef.DataList.Length = func() int {
		return len(*relatedData)
	}
	appRef.DataList.CreateItem = func() fyne.CanvasObject {
		return container.NewBorder(
			nil, nil, widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {}), nil,
			widget.NewButton("", func() {}))
	}
	// 重新定义List更新回调函数 然后刷新
	appRef.DataList.UpdateItem = func(i widget.ListItemID, o fyne.CanvasObject) {
		//Container中位置为0的元素是密码项按钮
		(o.(*fyne.Container).Objects[0]).(*widget.Button).SetText((*relatedData)[i].Name)
		(o.(*fyne.Container).Objects[0]).(*widget.Button).OnTapped = func() {
			// 根据回调函数提供的索引i，就是对应relatedData的当前点击项的顺序
			AppRef.ShowDataEditWinFunc(&(*relatedData)[i], cidOrg)
		}
		//Container中位置为1的元素是删除按钮
		(o.(*fyne.Container).Objects[1]).(*widget.Button).OnTapped = func() {
			dialog.ShowConfirm("提示", "确定要删除此条记录吗？", func(b bool) {
				if b {
					DeleteData((*relatedData)[i])
				}
			}, AppRef.W)
		}
	}
	appRef.DataList.Refresh()
}
