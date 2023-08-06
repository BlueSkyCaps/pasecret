package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"strconv"
)

func SendNotification(s string) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "Pasecret",
		Content: s,
	})
}

func main() {
	c := "12345"
	println(strconv.Atoi(string(c[0])))
	message01 := Test01("lockPwdSetTipShowConfirm", "zh")
	fmt.Println(message01)
}
