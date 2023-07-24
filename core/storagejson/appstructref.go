package storagejson

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
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
	dialog.ShowInformation("", "RepaintCartsByRemove", appRef.W)
	appRef.CardsGrid.Remove(delCard)
	appRef.CardsGrid.Refresh()
}

// RepaintCartsByEdit 成功保存本地存储库后再刷新Cart文件夹小部件
func (appRef AppStructRef) RepaintCartsByEdit(e *common.EditForm, editCard *widget.Card) {
	editCard.Title = e.Name
	editCard.Subtitle = e.Description
	dialog.ShowInformation("", "RepaintCartsByEdit", appRef.W)
	appRef.CardsGrid.Refresh()
}

// RepaintCartsByAdd 成功保存本地存储库后再刷新Cart文件夹小部件
func (appRef AppStructRef) RepaintCartsByAdd(addCi Category, addCard *widget.Card) {
	addCard.Title = addCi.Name
	addCard.Subtitle = addCi.Description
	dialog.ShowInformation("", "RepaintCartsByEdit", appRef.W)
	appRef.CardsGrid.Objects = append(appRef.CardsGrid.Objects, addCard)
	appRef.CardsGrid.Refresh()
}

// RepaintDataListByEdit 成功保存本地存储库后再刷新List密码项列表，避免本地保存失败却事先更新列表
func (appRef AppStructRef) RepaintDataListByEdit(performDataOrg *Data, cidOrg string) {
	// 获取此文件夹的密码项 最新数据
	relatedData := GetRelatedDataByCid(cidOrg)
	// 重新定义List更新回调函数 然后刷新
	appRef.DataList.UpdateItem = func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*widget.Button).SetText((*relatedData)[i].Name)
		o.(*widget.Button).OnTapped = func() {
			AppRef.ShowDataEditWinFunc(performDataOrg, cidOrg)
		}
	}
	appRef.DataList.Length = func() int {
		return len(*relatedData)
	}
	appRef.DataList.Refresh()
}
