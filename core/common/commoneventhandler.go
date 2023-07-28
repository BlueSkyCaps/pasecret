package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"sync"
	"time"
	"unicode/utf8"
)

func EntryOnChangedEventHandler(entry *widget.Entry) {
	if !fyne.CurrentDevice().IsMobile() {
		return
	}
	t := time.Now().UnixMilli()
	go func() {
		var tmp string
		var b sync.Mutex
		tmp = entry.Text
		entry.OnChanged = func(s string) {
			println("tmp:" + tmp)
			println("s:" + s)
			if utf8.RuneCountInString(s) >= utf8.RuneCountInString(tmp) {
				tmp = s
				return
			}
			b.Lock()
			if time.Now().UnixMilli()-t > 300 {
				//_, _ = utf8.DecodeLastRuneInString(tmp)
				entry.SetText(tmp)
				//tmp = entry.Text
				entry.CursorColumn = utf8.RuneCountInString(tmp)
				entry.Refresh()
			} else {
				tmp = entry.Text
			}
			t = time.Now().UnixMilli()
			b.Unlock()

		}
	}()
}
