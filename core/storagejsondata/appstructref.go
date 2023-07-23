package storagejsondata

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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

func (appRef AppStructRef) RepaintCarts(delCard *widget.Card) {
	dialog.ShowInformation("", "RepaintCarts", appRef.W)
	appRef.CardsGrid.Remove(delCard)
	appRef.CardsGrid.Refresh()
}
