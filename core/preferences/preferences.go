// Package preferences 加载首选项配置，如启动密码。不使用fyne的preferences提供的方法，因为在安卓端存储失败
package preferences

import (
	"encoding/json"
	"fyne.io/fyne/v2/dialog"
	"pasecret/core/common"
	"pasecret/core/storagejson"
	"path"
	"reflect"
)

type Preferences struct {
	LockPwd   string `json:"lock_pwd"`
	LocalLang string `json:"local_lang"`
}

var preferencePath string

func GetPreferenceByLockPwd() string {
	preferenceInit()
	preference := readPreference()
	return (*preference).LockPwd
}

func preferenceInit() {
	preferencePath = path.Join(storagejson.AppRef.A.Storage().RootURI().Path(), "preference.json")
	// 不存在首选项文件，则创建
	if !common.Existed(preferencePath) {
		initPre := Preferences{}
		marshal, err := json.Marshal(initPre)
		if err != nil {
			dialog.ShowInformation("err", "preferenceInit, json.Marshal:"+err.Error(), storagejson.AppRef.W)
			return
		}
		r, err := common.CreateFile(preferencePath, marshal)
		if !r {
			dialog.NewInformation("err", "preferenceInit, CreateFile:"+err.Error(), storagejson.AppRef.W).Show()
			return
		}
	}
}
func readPreference() *Preferences {
	r, bs, err := common.ReadFileAsBytes(preferencePath)
	if !r {
		dialog.NewInformation("err", "readPreference,ReadFileAsBytes:"+err.Error(), storagejson.AppRef.W).Show()
		return nil
	}
	preference := Preferences{}
	err = json.Unmarshal(bs, &preference)
	if err != nil {
		dialog.NewInformation("err", "readPreference, json.Marshal d:"+err.Error(), storagejson.AppRef.W).Show()
		return nil
	}
	return &preference
}
func SetPreference(key string, v interface{}) {
	preferenceInit()
	preference := readPreference()
	preferenceR := reflect.ValueOf(preference)
	keyNameR := preferenceR.Elem().FieldByName(key)
	if !reflect.ValueOf(v).Type().AssignableTo(keyNameR.Type()) {
		dialog.NewInformation("err", "AssignableTo, v ref cant assignable to keyNameR", storagejson.AppRef.W).Show()
		return
	}
	keyNameR.Set(reflect.ValueOf(v))
	marshal, err := json.Marshal(preference)
	if err != nil {
		dialog.ShowInformation("err", "SetPreferenceByLockPwd, json.Marshal:"+err.Error(), storagejson.AppRef.W)
		return
	}
	r, err := common.WriteExistedFile(preferencePath, marshal)
	if !r {
		dialog.NewInformation("err", "SetPreferenceByLockPwd, WriteExistedFile:"+err.Error(), storagejson.AppRef.W).Show()
		return
	}
}
func RemovePreference(key string) {
	preference := readPreference()
	preferenceR := reflect.ValueOf(preference)
	keyNameR := preferenceR.Elem().FieldByName(key)
	if keyNameR.IsZero() {
		dialog.NewInformation("err", "RemovePreferenceBy, key ref is zero", storagejson.AppRef.W).Show()
		return
	}
	keyNameR.Set(reflect.Zero(keyNameR.Type()))
	marshal, err := json.Marshal(preference)
	if err != nil {
		dialog.ShowInformation("err", "RemovePreferenceByLockPwd, json.Marshal:"+err.Error(), storagejson.AppRef.W)
		return
	}
	r, err := common.WriteExistedFile(preferencePath, marshal)
	if !r {
		dialog.NewInformation("err", "RemovePreferenceByLockPwd, WriteExistedFile:"+err.Error(), storagejson.AppRef.W).Show()
		return
	}
}
