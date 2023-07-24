package storagejson

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
)

// AppStructRef App操作句柄
type AppStructRef struct {
	W           fyne.Window
	A           fyne.App
	LoadedItems LoadedItems
	SearchInput *widget.Entry
	SearchBtn   *widget.Button
	CardsGrid   *fyne.Container
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
