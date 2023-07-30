// Package preferences 加载首选项配置，如启动密码。不使用fyne的preferences提供的方法，因为在安卓端存储失败
package preferences

import (
	"encoding/json"
	"fyne.io/fyne/v2/dialog"
	"pasecret/core/common"
	"pasecret/core/storagejson"
	"path"
)

type Preferences struct {
	LockPwd string `json:"lock_pwd"`
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
		initPre := Preferences{
			LockPwd: "",
		}
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
func SetPreferenceByLockPwd(v interface{}) {
	preferenceInit()
	preference := readPreference()
	preference.LockPwd = v.(string)
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
func RemovePreferenceByLockPwd() {
	preferenceInit()
	preference := readPreference()
	preference.LockPwd = ""
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
