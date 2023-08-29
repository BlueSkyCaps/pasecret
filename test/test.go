package main

import (
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
	c := make(chan int, 2)
	go func() {
		time.Sleep(2 * time.Second)
		c <- 1
		defer println("111")
	}()
	go func() {
		time.Sleep(4 * time.Second)
		c <- 1
		close(c)
		defer println("111")
	}()
	<-c
	time.Sleep(6 * time.Second)
	println("222")

	<-c
	println("222")

	//
	//println(common.EncryptAES([]byte(common.AppProductKeyAES), "1234"))
	//println(common.EncryptAES([]byte(common.AppProductKeyAES), "1234"))
	//println(common.DecryptAES([]byte(common.AppProductKeyAES), "vHg+/9NnScW6M1nYdm/MyU5loM4="))
	//println(common.DecryptAES([]byte(common.AppProductKeyAES), "IZGyzjZf+/RcLnx6pF/vW5ezlEA=1"))

	//c := "12345"
	//println(strconv.Atoi(string(c[0])))
	//message01 := Test01("lockPwdSetTipShowConfirm", "zh")
	//fmt.Println(message01)
}
