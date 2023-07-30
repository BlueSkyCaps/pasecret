package test

import (
	"fmt"
	"fyne.io/fyne/v2"
	"time"
)

func SendNotification(s string) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "Pasecret",
		Content: s,
	})
}

func main() {
	p := time.Now().UnixMilli()
	time.Sleep(time.Millisecond * 3000)
	fmt.Println(time.Now().UnixMilli())
	fmt.Println(p)

	s := []string{"1", "2", "3"}
	fmt.Println(s[:len(s)-1])
}
