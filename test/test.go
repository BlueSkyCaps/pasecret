package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"pasecret/core/i18n"
)

func SendNotification(s string) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "Pasecret",
		Content: s,
	})
}

func main() {
	message01 := i18n.Test01("Tips", "zh")
	fmt.Println(message01)
}
