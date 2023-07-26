// Package storagejson 贮存密码箱数据，加载、保存本地数据
package storagejson

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"path"
	"sort"
	"strings"
	"time"
)

var AppRef AppStructRef
var stoDPath string

func LoadInit(data []byte) {
	stoDPath = path.Join(AppRef.A.Storage().RootURI().Path(), "d.json")
	existed := common.Existed(AppRef.A.Storage().RootURI().Path())
	// 若RootURI目录不存在，则先创建目录（目前是在Android端必须先创建，因为不存在/data/user/0/top.reminisce.xxx/files/fyne）
	if !existed {
		r, err := common.CreateDir(AppRef.A.Storage().RootURI().Path())
		if !r {
			dialog.NewInformation("err", "storage loadInit, MkdirAll:"+err.Error(), AppRef.W).Show()
		}
	}

	stoDPath = path.Join(AppRef.A.Storage().RootURI().Path(), "d.json")
	// 不存在默认密码数据文件，则创建
	if !common.Existed(stoDPath) {
		r, err := common.CreateFile(stoDPath, data)
		if !r {
			dialog.NewInformation("err", "storage loadInit, CreateFile:"+err.Error(), AppRef.W).Show()
		}
	}
	load(stoDPath)
}

// 从本地存储库d.json加载密码数据
func load(stoDPath string) {
	r, bs, err := common.ReadFileAsBytes(stoDPath)
	if !r {
		dialog.NewInformation("err", "storage load, ReadFileAsString:"+err.Error(), AppRef.W).Show()
	}
	err = json.Unmarshal(bs, &AppRef.LoadedItems)
	if err != nil {
		dialog.NewInformation("err", "storage load, json.Marshal d:"+err.Error(), AppRef.W).Show()
	}
	// 解密
	decLoadedData()
}

// 解密所有密码项重新贮存AppRef
func decLoadedData() {
	var err error

	for i := 0; i < len(AppRef.LoadedItems.Data); i++ {
		AppRef.LoadedItems.Data[i].Name, err = common.DecryptAES([]byte(common.AppProductKeyAES), AppRef.LoadedItems.Data[i].Name)
		AppRef.LoadedItems.Data[i].Site, err = common.DecryptAES([]byte(common.AppProductKeyAES), AppRef.LoadedItems.Data[i].Site)
		AppRef.LoadedItems.Data[i].Remark, err = common.DecryptAES([]byte(common.AppProductKeyAES), AppRef.LoadedItems.Data[i].Remark)
		AppRef.LoadedItems.Data[i].Password, err = common.DecryptAES([]byte(common.AppProductKeyAES), AppRef.LoadedItems.Data[i].Password)
		AppRef.LoadedItems.Data[i].AccountName, err = common.DecryptAES([]byte(common.AppProductKeyAES), AppRef.LoadedItems.Data[i].AccountName)
	}
	if err != nil {
		dialog.ShowInformation("终止", "decLoadedData, common.DecryptAES:"+err.Error(), AppRef.W)
		time.Sleep(time.Millisecond * 5000)
		go func() {
			time.Sleep(time.Millisecond * 4800)
			AppRef.W.Close()
		}()
	}
}

// 加密所有密码项再保存至本地存储库
func encryLoadedData() {
	var err error
	// 加密所有密码项
	for i := 0; i < len(AppRef.LoadedItems.Data); i++ {
		AppRef.LoadedItems.Data[i].Name, err = common.EncryptAES([]byte(common.AppProductKeyAES), AppRef.LoadedItems.Data[i].Name)
		AppRef.LoadedItems.Data[i].Site, err = common.EncryptAES([]byte(common.AppProductKeyAES), AppRef.LoadedItems.Data[i].Site)
		AppRef.LoadedItems.Data[i].Remark, err = common.EncryptAES([]byte(common.AppProductKeyAES), AppRef.LoadedItems.Data[i].Remark)
		AppRef.LoadedItems.Data[i].Password, err = common.EncryptAES([]byte(common.AppProductKeyAES), AppRef.LoadedItems.Data[i].Password)
		AppRef.LoadedItems.Data[i].AccountName, err = common.EncryptAES([]byte(common.AppProductKeyAES), AppRef.LoadedItems.Data[i].AccountName)
	}
	if err != nil {
		dialog.ShowInformation("终止", "encryLoadedData, common.EncryptAES:"+err.Error(), AppRef.W)
		time.Sleep(time.Millisecond * 5000)
		go func() {
			time.Sleep(time.Millisecond * 4800)
			AppRef.W.Close()
		}()
	}
}

// EditCategory 编辑保存一个归类文件夹
func EditCategory(e *common.EditForm, editCi Category, editCard *widget.Card) {
	var newCategory []Category
	for _, ci := range AppRef.LoadedItems.Category {
		if ci.Id == editCi.Id {
			ci.Name = e.Name
			ci.Alias = e.Alias
			ci.Description = e.Description
		}
		newCategory = append(newCategory, ci)
	}
	AppRef.LoadedItems.Category = newCategory
	// 先加密保存
	encryLoadedData()
	marshalDJson, err := json.Marshal(AppRef.LoadedItems)
	if err != nil {
		dialog.NewInformation("err", "EditCategory, json.Marshal:"+err.Error(), AppRef.W).Show()
		return
	}
	r, err := common.WriteExistedFile(stoDPath, marshalDJson)
	if !r {
		dialog.NewInformation("err", "EditCategory, WriteExistedFile:"+err.Error(), AppRef.W).Show()
		return
	}
	// 保存完后要解密重新贮存至AppRef
	decLoadedData()
	// 成功保存本地存储库后再刷新Cart文件夹小部件，避免本地保存失败却事先更新Card小部件文本
	AppRef.RepaintCartsByEdit(e, editCard)
	return
}

// AddCategory 新增一个归类文件夹，但此函数不更新窗口此Cart小部件
func AddCategory(e *common.EditForm) *Category {
	var addCategory Category
	addCategory.Name = e.Name
	addCategory.Alias = e.Alias
	addCategory.Description = e.Description
	addCategory.Renameable = true
	addCategory.Removable = true
	r, rank := common.GenAscRankId()
	if !r {
		dialog.NewInformation("err", "AddCategory, GenAscRankId", AppRef.W).Show()
		return nil
	}
	addCategory.Rank = rank
	addCategory.Id = fmt.Sprintf("%d-built-in-can-removed", addCategory.Rank)
	AppRef.LoadedItems.Category = append(AppRef.LoadedItems.Category, addCategory)
	encryLoadedData()
	marshalDJson, err := json.Marshal(AppRef.LoadedItems)
	if err != nil {
		dialog.NewInformation("err", "AddCategory, json.Marshal:"+err.Error(), AppRef.W).Show()
		return nil
	}
	r, err = common.WriteExistedFile(stoDPath, marshalDJson)
	if !r {
		dialog.NewInformation("err", "AddCategory, WriteExistedFile:"+err.Error(), AppRef.W).Show()
		return nil
	}
	// 保存完后要解密重新贮存至AppRef
	decLoadedData()
	return &addCategory
}

// DeleteCategory 删除一个归类文件夹
func DeleteCategory(delCi Category, delCard *widget.Card) {
	var newCategory []Category
	for _, ci := range AppRef.LoadedItems.Category {
		if ci.Id != delCi.Id {
			newCategory = append(newCategory, ci)
		}
	}
	AppRef.LoadedItems.Category = newCategory
	deleteCategoryRelated(delCi)
	encryLoadedData()
	marshalDJson, err := json.Marshal(AppRef.LoadedItems)
	if err != nil {
		dialog.NewInformation("err", "DeleteCategory, json.Marshal:"+err.Error(), AppRef.W).Show()
		return
	}
	r, err := common.WriteExistedFile(stoDPath, marshalDJson)
	if !r {
		dialog.NewInformation("err", "DeleteCategory, WriteExistedFile:"+err.Error(), AppRef.W).Show()
		return
	}
	// 保存完后要解密重新贮存至AppRef
	decLoadedData()
	AppRef.RepaintCartsByRemove(delCard)
	return
}

// 删除一个归类文件夹下的所有密码项
func deleteCategoryRelated(delCi Category) {
	var newData []Data
	for _, da := range AppRef.LoadedItems.Data {
		if da.CategoryId != delCi.Id {
			newData = append(newData, da)
		}
	}
	AppRef.LoadedItems.Data = newData
}

// SearchByKeywordTo 按关键词搜索密码项，返回新的LoadedItems副本供部件展示
func SearchByKeywordTo(kw string) []common.SearchDataResultViewModel {
	var vm []common.SearchDataResultViewModel
	for _, d := range AppRef.LoadedItems.Data {
		if strings.Contains(d.Name, kw) || strings.Contains(d.Remark, kw) || strings.Contains(d.AccountName, kw) ||
			strings.Contains(d.Site, kw) {
			var c common.SearchDataResultViewModel
			c.VDataAccountName = d.AccountName
			// 根据CategoryId找归类文件夹名称
			for _, ci := range AppRef.LoadedItems.Category {
				if ci.Id == d.CategoryId {
					c.VDataCategoryName = ci.Name
				}
			}
			c.VDataName = d.Name
			c.VDataPassword = d.Password
			c.VDataRemark = d.Remark
			c.VDataSite = d.Site
			vm = append(vm, c)
		}
	}
	return vm
}

// 按原有Rank顺序排序
func sortCategory() {
	sort.Slice(AppRef.LoadedItems.Category, func(i, j int) bool {
		return AppRef.LoadedItems.Category[i].Rank < AppRef.LoadedItems.Category[i].Rank
	})
}
