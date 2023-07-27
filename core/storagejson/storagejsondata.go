package storagejson

import (
	"encoding/json"
	"fyne.io/fyne/v2/dialog"
	"pasecret/core/common"
)

// GetRelatedDataByCid 根据归类文件夹Id检索它的密码项
func GetRelatedDataByCid(cid string) *[]Data {
	// new返回指向[]Data的指针
	var relatedData = new([]Data)
	for _, d := range AppRef.LoadedItems.Data {
		if d.CategoryId == cid {
			*relatedData = append(*relatedData, d)
		}
	}
	return relatedData
}

// EditData 编辑保存一个密码项，包括新增
func EditData(theData Data, isEditOp bool, cidOrg string) {
	if isEditOp {
		var newData []Data
		// 编辑则更新AppRef中的此密码项
		for _, d := range AppRef.LoadedItems.Data {
			if d.Id == theData.Id {
				d.Name = theData.Name
				d.Site = theData.Site
				d.Password = theData.Password
				d.AccountName = theData.AccountName
				d.Remark = theData.Remark
			}
			newData = append(newData, d)
		}
		AppRef.LoadedItems.Data = newData

	} else {
		// 新增则往里面追加
		AppRef.LoadedItems.Data = append(AppRef.LoadedItems.Data, theData)
	}
	encryLoadedData()
	marshalDJson, err := json.Marshal(AppRef.LoadedItems)
	if err != nil {
		dialog.NewInformation("err", "EditCategory, json.Marshal:"+err.Error(), AppRef.W).Show()
		return
	}
	r, err := common.WriteExistedFile(StoDPath, marshalDJson)
	if !r {
		dialog.NewInformation("err", "EditCategory, WriteExistedFile:"+err.Error(), AppRef.W).Show()
		return
	}
	// 保存完后要解密重新贮存至AppRef
	decLoadedData()
	// 成功保存本地存储库后再刷新List密码项列表，避免本地保存失败却事先更新列表
	AppRef.RepaintDataListByEdit(cidOrg)
	return
}

// DeleteData 删除一个密码项
func DeleteData(delD Data) {
	var newData []Data
	for _, di := range AppRef.LoadedItems.Data {
		if di.Id != delD.Id {
			newData = append(newData, di)
		}
	}
	AppRef.LoadedItems.Data = newData
	encryLoadedData()
	marshalDJson, err := json.Marshal(AppRef.LoadedItems)
	if err != nil {
		dialog.NewInformation("err", "DeleteData, json.Marshal:"+err.Error(), AppRef.W).Show()
		return
	}
	r, err := common.WriteExistedFile(StoDPath, marshalDJson)
	if !r {
		dialog.NewInformation("err", "DeleteData, WriteExistedFile:"+err.Error(), AppRef.W).Show()
		return
	}
	// 保存完后要解密重新贮存至AppRef
	decLoadedData()
	AppRef.RepaintDataListByEdit(delD.CategoryId)
	return
}
