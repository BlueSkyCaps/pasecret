// Package storagedata 贮存密码箱数据，加载、保存本地存储库
package storagedata

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"path"
	"strconv"
	"strings"
	"time"
)

var AppRef AppStructRef
var StoDPath string

// LoadInit 初始化載入本地存储库數據文件
func LoadInit(data []byte) {
	StoDPath = path.Join(AppRef.A.Storage().RootURI().Path(), "d.json")
	existed := common.Existed(AppRef.A.Storage().RootURI().Path())
	// 若RootURI目录不存在，则先创建目录（目前是在Android端必须先创建，因为不存在/data/user/0/top.reminisce.xxx/files/fyne）
	if !existed {
		r, err := common.CreateDir(AppRef.A.Storage().RootURI().Path())
		if !r {
			dialog.NewInformation("err", "storage loadInit, MkdirAll:"+err.Error(), AppRef.W).Show()
			return
		}
	}

	StoDPath = path.Join(AppRef.A.Storage().RootURI().Path(), "d.json")
	// 不存在默认数据文件，则创建
	if !common.Existed(StoDPath) {
		r, err := common.CreateFile(StoDPath, data)
		if !r {
			dialog.NewInformation("err", "storage loadInit, CreateFile:"+err.Error(), AppRef.W).Show()
			return
		}
	}
	load(StoDPath)
}

// 从本地存储库d.json加载密码数据
func load(stoDPath string) {
	r, bs, err := common.ReadFileAsBytes(stoDPath)
	if !r {
		dialog.NewInformation("err", "storage load, ReadFileAsString:"+err.Error(), AppRef.W).Show()
		return
	}
	err = json.Unmarshal(bs, &AppRef.LoadedItems)
	if err != nil {
		dialog.NewInformation("err", "storage load, json.Marshal d:"+err.Error(), AppRef.W).Show()
		return
	}
	// 解密密码项
	decLoadedData()
	// 根据语言环境设置内置归类文件夹显示的文本
	setCategoryToCartTextByLang()
}

// 若是内置归类夹，不能删除和编辑。根据当前语言环境设置标题和描述
func setCategoryToCartTextByLang() {
	var err error
	for i := 0; i < len(AppRef.LoadedItems.Category); i++ {
		if !strings.HasSuffix(AppRef.LoadedItems.Category[i].Id, "built-in-cannot-be-removed") {
			continue
		}
		c := 0
		c, err = strconv.Atoi(string(AppRef.LoadedItems.Category[i].Id[0]))
		// 变量c可以区分国际化文本的message id
		AppRef.LoadedItems.Category[i].Name = AppRef.LocalizedTextFunc(fmt.Sprintf("buildInCategory%dName", c), nil)
		AppRef.LoadedItems.Category[i].Description = AppRef.LocalizedTextFunc(fmt.Sprintf("buildInCategory%dDescription", c), nil)
	}
	if err != nil {
		dialog.ShowInformation("err", "setCategoryToCartTextByLang:"+err.Error(), AppRef.W)
		time.Sleep(time.Millisecond * 5000)
		go func() {
			time.Sleep(time.Millisecond * 4800)
			AppRef.W.Close()
		}()
	}
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
		dialog.ShowInformation("err", "decLoadedData, common.DecryptAES:"+err.Error(), AppRef.W)
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
		dialog.ShowInformation("err", "encryLoadedData, common.EncryptAES:"+err.Error(), AppRef.W)
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
	r, err := common.WriteExistedFile(StoDPath, marshalDJson)
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
	r, err = common.WriteExistedFile(StoDPath, marshalDJson)
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
	r, err := common.WriteExistedFile(StoDPath, marshalDJson)
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
		if strings.Contains(strings.ToLower(d.Name), kw) || strings.Contains(strings.ToLower(d.Remark), kw) ||
			strings.Contains(strings.ToLower(d.AccountName), kw) || strings.Contains(strings.ToLower(d.Site), kw) {
			var c common.SearchDataResultViewModel
			c.VDataAccountName = d.AccountName
			// 根据CategoryId找归类文件夹名称
			realCi := GetCategoryByCid(d.CategoryId)
			c.VDataCategoryName = realCi.Name
			c.VDataName = d.Name
			c.VDataPassword = d.Password
			c.VDataRemark = d.Remark
			c.VDataSite = d.Site
			vm = append(vm, c)
		}
	}
	return vm
}

func GetCategoryByCid(cid string) Category {
	var realCi Category
	for _, nci := range AppRef.LoadedItems.Category {
		if nci.Id == cid {
			realCi = nci
		}
	}
	return realCi
}
