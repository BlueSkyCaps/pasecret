package storagejsondata

import (
	"fyne.io/fyne/v2"
)

// AppStructRef App操作句柄
type AppStructRef struct {
	W           fyne.Window
	A           fyne.App
	LoadedItems LoadedItems
}
