// Package preferences 加载首选项配置，如启动密码。不使用fyne的preferences提供的方法，因为在安卓端存储失败
package preferences

import (
	"encoding/json"
	"fyne.io/fyne/v2/dialog"
	"pasecret/core/common"
	"pasecret/core/storagedata"
	"path"
	"reflect"
)

type Preferences struct {
	LockPwd   string `json:"lock_pwd"`
	LocalLang string `json:"local_lang"`
}

var preferencePath string

func GetPreferenceByLockPwd() string {
	PreferenceInit()
	preference := readPreference()
	return (*preference).LockPwd
}

func GetPreferenceByLocalLang() string {
	PreferenceInit()
	preference := readPreference()
	return (*preference).LocalLang
}

func PreferenceInit() {
	existed := common.Existed(storagedata.AppRef.A.Storage().RootURI().Path())
	// 若RootURI目录不存在，则先创建目录（目前是在Android端必须先创建，因为不存在/data/user/0/top.reminisce.xxx/files/fyne）
	if !existed {
		r, err := common.CreateDir(storagedata.AppRef.A.Storage().RootURI().Path())
		if !r {
			dialog.NewInformation("err", "storage loadInit, MkdirAll:"+err.Error(), storagedata.AppRef.W).Show()
			return
		}
	}
	preferencePath = path.Join(storagedata.AppRef.A.Storage().RootURI().Path(), "preference.json")

	// 不存在首选项文件，则创建
	if !common.Existed(preferencePath) {
		initPre := Preferences{}
		marshal, err := json.Marshal(initPre)
		if err != nil {
			dialog.ShowInformation("err", "PreferenceInit, json.Marshal:"+err.Error(), storagedata.AppRef.W)
			return
		}
		r, err := common.CreateFile(preferencePath, marshal)
		if !r {
			dialog.NewInformation("err", "PreferenceInit, CreateFile:"+err.Error(), storagedata.AppRef.W).Show()
			return
		}
	}
}
func readPreference() *Preferences {
	r, bs, err := common.ReadFileAsBytes(preferencePath)
	if !r {
		dialog.NewInformation("err", "readPreference,ReadFileAsBytes:"+err.Error(), storagedata.AppRef.W).Show()
		return nil
	}
	preference := Preferences{}
	err = json.Unmarshal(bs, &preference)
	if err != nil {
		dialog.NewInformation("err", "readPreference, json.Marshal d:"+err.Error(), storagedata.AppRef.W).Show()
		return nil
	}
	return &preference
}

// SetPreference 设置首选项Preferences某个键的值，注意fieldName是结构体字段名而不是json文件实际存储的key键名
func SetPreference(fieldName string, v interface{}) {
	PreferenceInit()
	preference := readPreference()
	preferenceR := reflect.ValueOf(preference)
	keyNameR := preferenceR.Elem().FieldByName(fieldName)
	if !reflect.ValueOf(v).Type().AssignableTo(keyNameR.Type()) {
		dialog.NewInformation("err", "AssignableTo, v ref cant assignable to keyNameR", storagedata.AppRef.W).Show()
		return
	}
	keyNameR.Set(reflect.ValueOf(v))
	marshal, err := json.Marshal(preference)
	if err != nil {
		dialog.ShowInformation("err", "SetPreferenceByLockPwd, json.Marshal:"+err.Error(), storagedata.AppRef.W)
		return
	}
	r, err := common.WriteExistedFile(preferencePath, marshal)
	if !r {
		dialog.NewInformation("err", "SetPreferenceByLockPwd, WriteExistedFile:"+err.Error(), storagedata.AppRef.W).Show()
		return
	}
}
func RemovePreference(key string) {
	preference := readPreference()
	preferenceR := reflect.ValueOf(preference)
	keyNameR := preferenceR.Elem().FieldByName(key)
	if keyNameR.IsZero() {
		dialog.NewInformation("err", "RemovePreferenceBy, key ref is zero", storagedata.AppRef.W).Show()
		return
	}
	keyNameR.Set(reflect.Zero(keyNameR.Type()))
	marshal, err := json.Marshal(preference)
	if err != nil {
		dialog.ShowInformation("err", "RemovePreferenceByLockPwd, json.Marshal:"+err.Error(), storagedata.AppRef.W)
		return
	}
	r, err := common.WriteExistedFile(preferencePath, marshal)
	if !r {
		dialog.NewInformation("err", "RemovePreferenceByLockPwd, WriteExistedFile:"+err.Error(), storagedata.AppRef.W).Show()
		return
	}
}
