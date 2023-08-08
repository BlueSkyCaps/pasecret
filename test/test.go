package main

import (
	"fyne.io/fyne/v2"
	"pasecret/core/common"
)

func SendNotification(s string) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "Pasecret",
		Content: s,
	})
}

func main() {
	println(common.EncryptAES([]byte(common.AppProductKeyAES), "1234"))
	println(common.EncryptAES([]byte(common.AppProductKeyAES), "1234"))
	println(common.DecryptAES([]byte(common.AppProductKeyAES), "vHg+/9NnScW6M1nYdm/MyU5loM4="))
	println(common.DecryptAES([]byte(common.AppProductKeyAES), "IZGyzjZf+/RcLnx6pF/vW5ezlEA=1"))

	//c := "12345"
	//println(strconv.Atoi(string(c[0])))
	//message01 := Test01("lockPwdSetTipShowConfirm", "zh")
	//fmt.Println(message01)
}
