// Package storagejsondata 贮存密码箱数据，加载、保存本地数据
package storagejsondata

import (
	"encoding/json"
	"fyne.io/fyne/v2/dialog"
	"pasecret/core/common"
	"path"
)

func LoadInit(appRef AppStructRef, data []byte) LoadedItems {
	existed := common.Existed(appRef.A.Storage().RootURI().Path())
	// 若RootURI目录不存在，则先创建目录（目前是在Android端必须先创建，因为不存在/data/user/0/top.reminisce.xxx/files/fyne）
	if !existed {
		r, err := common.CreateDir(appRef.A.Storage().RootURI().Path())
		if !r {
			dialog.NewInformation("err", "storage loadInit, MkdirAll:"+err.Error(), appRef.W).Show()
		}
	}

	stoDPath := path.Join(appRef.A.Storage().RootURI().Path(), "d.json")
	// 不存在默认密码数据文件，则创建
	if !common.Existed(stoDPath) {
		r, err := common.CreateFile(stoDPath, data)
		if !r {
			dialog.NewInformation("err", "storage loadInit, CreateFile:"+err.Error(), appRef.W).Show()
		}
	}
	return load(stoDPath, appRef)
}

// 从本地存储库d.json加载密码数据
func load(stoDPath string, appRef AppStructRef) LoadedItems {
	r, bs, err := common.ReadFileAsBytes(stoDPath)
	if !r {
		dialog.NewInformation("err", "storage load, ReadFileAsString:"+err.Error(), appRef.W).Show()
	}
	var loadedItems LoadedItems
	err = json.Unmarshal(bs, &loadedItems)
	if err != nil {
		dialog.NewInformation("err", "storage load, json.Marshal d:"+err.Error(), appRef.W).Show()
	}
	println(bs)
	dialog.ShowInformation("", appRef.A.Storage().RootURI().Path(), appRef.W)
	return loadedItems
}
