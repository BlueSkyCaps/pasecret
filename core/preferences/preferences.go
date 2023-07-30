// Package preferences 加载首选项配置，如启动密码
package preferences

import "pasecret/core/storagejson"

func GetPreferenceByLockPwd() string {
	return storagejson.AppRef.A.Preferences().String("lock_pwd")
}

func SetPreferenceByLockPwd(v interface{}) {
	storagejson.AppRef.A.Preferences().SetString("lock_pwd", v.(string))
}
