// Package storagejsondata 贮存密码箱数据，加载、保存本地数据
package storagejsondata

import (
	"encoding/json"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"path"
)

var AppRef AppStructRef
var stoDPath string
var loadedItems LoadedItems

func LoadInit(appRef AppStructRef, data []byte) LoadedItems {
	AppRef = appRef
	stoDPath = path.Join(AppRef.A.Storage().RootURI().Path(), "d.json")
	existed := common.Existed(AppRef.A.Storage().RootURI().Path())
	// 若RootURI目录不存在，则先创建目录（目前是在Android端必须先创建，因为不存在/data/user/0/top.reminisce.xxx/files/fyne）
	if !existed {
		r, err := common.CreateDir(AppRef.A.Storage().RootURI().Path())
		if !r {
			dialog.NewInformation("err", "storage loadInit, MkdirAll:"+err.Error(), AppRef.W).Show()
		}
	}

	stoDPath := path.Join(AppRef.A.Storage().RootURI().Path(), "d.json")
	// 不存在默认密码数据文件，则创建
	if !common.Existed(stoDPath) {
		r, err := common.CreateFile(stoDPath, data)
		if !r {
			dialog.NewInformation("err", "storage loadInit, CreateFile:"+err.Error(), AppRef.W).Show()
		}
	}
	return load(stoDPath, AppRef)
}

// 从本地存储库d.json加载密码数据
func load(stoDPath string, AppRef AppStructRef) LoadedItems {
	r, bs, err := common.ReadFileAsBytes(stoDPath)
	if !r {
		dialog.NewInformation("err", "storage load, ReadFileAsString:"+err.Error(), AppRef.W).Show()
	}
	err = json.Unmarshal(bs, &loadedItems)
	if err != nil {
		dialog.NewInformation("err", "storage load, json.Marshal d:"+err.Error(), AppRef.W).Show()
	}
	println(bs)
	dialog.ShowInformation("", AppRef.A.Storage().RootURI().Path(), AppRef.W)
	return loadedItems
}

// DeleteCategory 删除一个归类文件夹
func DeleteCategory(delCi Category, delCard *widget.Card) {
	var newCategory []Category
	for _, ci := range loadedItems.Category {
		if ci.Id != delCi.Id {
			newCategory = append(newCategory, ci)
		}
	}
	loadedItems.Category = newCategory
	deleteCategoryRelated(delCi)
	marshalDJson, err := json.Marshal(loadedItems)
	if err != nil {
		dialog.NewInformation("err", "DeleteCategory, json.Marshal:"+err.Error(), AppRef.W).Show()
		return
	}
	r, err := common.WriteExistedFile(stoDPath, marshalDJson)
	if !r {
		dialog.NewInformation("err", "DeleteCategory, WriteExistedFile:"+err.Error(), AppRef.W).Show()
		return
	}
	AppRef.RepaintCarts(delCard)
	return
}

// 删除一个归类文件夹下的所有密码项
func deleteCategoryRelated(delCi Category) {
	var newData []Data
	for _, da := range loadedItems.Data {
		if da.CategoryId != delCi.Id {
			newData = append(newData, da)
		}
	}
	loadedItems.Data = newData
}
