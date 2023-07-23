package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"io/ioutil"
	"os"
	"pasecret/core/ui"
	"path"
)

var w fyne.Window
var a fyne.App

func testStorage() {

	_, err := os.Stat(a.Storage().RootURI().Path())
	// 若RootURI目录不存在，则先创建目录（目前是在Android端必须先创建，因为不存在/data/user/0/top.reminisce.xxx/files/fyne）
	if err != nil {
		err = os.MkdirAll(a.Storage().RootURI().Path(), os.ModePerm)
		if err != nil {
			dialog.NewInformation("", "MkdirAll:"+err.Error(), w).Show()
		}
	}

	stoPath := path.Join(a.Storage().RootURI().Path(), "aa2.txt")
	// Android必须添加os.O_WRONLY，否则报"bad file descriptor"权限问题
	create, err := os.OpenFile(stoPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		dialog.NewInformation("", "OpenFile:"+err.Error(), w).Show()
	}
	_, err = create.Write([]byte("test 1307我爱你2"))
	if err != nil {
		dialog.NewInformation("", "Write:"+err.Error(), w).Show()
	}
	err = create.Close()
	if err != nil {
		dialog.NewInformation("", "create.Close():"+err.Error(), w).Show()
	}
	b, err := ioutil.ReadFile(stoPath)
	if err != nil {
		dialog.NewInformation("", "open.Read(ra):"+err.Error(), w).Show()

	}

	v := string(b)
	dialog.ShowInformation("", string(v), w)
	dialog.ShowInformation("", a.Storage().RootURI().Path(), w)
	ui.Run()

}
