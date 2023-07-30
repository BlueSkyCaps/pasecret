package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"sync"
	"time"
	"unicode/utf8"
)

// EntryOnChangedEventHandler 安卓端退格键中文字符会被执行两次，导致删除两个字符（MuMu模拟机器中所有字符都会被退格两次）
// 也许是fyne存在的bug。 Windows不会有此问题。
//此Handler监听传递的Entry文本改变事件，根据毫秒级时间戳判断触发间隔来设置文本回退避免被删除两次
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
			if utf8.RuneCountInString(s) >= utf8.RuneCountInString(tmp) {
				tmp = s
				return
			}
			b.Lock()
			if time.Now().UnixMilli()-t > 200 {
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

var SearchTmp string

func SearchEntryOnChangedEventHandler(entry *widget.Entry) {
	if !fyne.CurrentDevice().IsMobile() {
		return
	}
	t := time.Now().UnixMilli()
	var b sync.Mutex
	SearchTmp = entry.Text
	entry.OnChanged = func(s string) {
		if utf8.RuneCountInString(s) >= utf8.RuneCountInString(SearchTmp) {
			SearchTmp = s
			return
		}
		b.Lock()
		if time.Now().UnixMilli()-t > 200 {
			//_, _ = utf8.DecodeLastRuneInString(tmp)
			entry.SetText(SearchTmp)
			//tmp = entry.Text
			entry.CursorColumn = utf8.RuneCountInString(SearchTmp)
			entry.Refresh()
		} else {
			SearchTmp = entry.Text
		}
		t = time.Now().UnixMilli()
		b.Unlock()
	}
}
