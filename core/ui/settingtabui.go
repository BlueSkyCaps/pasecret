package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"pasecret/core/storagejson"
	"strings"
	"time"
)

var treeSettingDictParent = []string{"启动密码", "备份还原", "多端同步", "捐助"}
var treeSettingDict map[string][]string = map[string][]string{
	treeSettingDictParent[0]: {"NotChild_0"},
	treeSettingDictParent[1]: {"备份数据", "还原数据到本机"},
	treeSettingDictParent[2]: {"NotChild_1"},
	treeSettingDictParent[3]: {"NotChild_3"},
}

func createSettingTabContent() *widget.Tree {
	tree := widget.NewTree(
		// 定义所有分支的唯一id，来自treeSettingDictParent和treeSettingDict
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			if id == "" {
				return treeSettingDictParent
			}
			return treeSettingDict[id]
		},
		// 定义是否有子元素或者没有
		func(id widget.TreeNodeID) bool {
			// ""空字符串约定是tree最外层id，不可忽略
			if id == "" {
				return true
			}
			/*treeSettingDictParent的元素都是父节点它们可能有子节点
			规定treeSettingDict[pid][0]若是NotChild_表示没有子元素节点*/
			for _, pid := range treeSettingDictParent {
				if id == pid && !strings.HasPrefix(treeSettingDict[pid][0], "NotChild_") {
					return true
				}
			}
			// 其余的id都是最底层元素，没有子节点
			return false
		},
		// 定义树节点的小部件
		func(branch bool) fyne.CanvasObject {
			if branch {
				return widget.NewButton("Branch template", func() {})
			}
			return widget.NewButton("Leaf template", func() {})
		},
		// 真正更新树节点的小部件，定义点击响应事件
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			button := o.(*widget.Button)
			button.SetText(id)
			// 若是备份数据按钮
			if id == treeSettingDict[treeSettingDictParent[1]][0] {
				button.OnTapped = backupBthCallBack
			}
			// 若是还原数据按钮
			if id == treeSettingDict[treeSettingDictParent[1]][1] {
				button.OnTapped = restoreBthCallBack
			}

		})
	return tree
}

// 备份数据响应函数
/*
	安卓端必须使用fyne.URIWriteCloser.Write将数据写入选择的uri，或许是fyne的问题，
	安卓端无法将数据写入到非app本地存储库中，common.CreateFile函数创建uriWriteCloser.URI().Path()的目录提示目录不存在或权限问题。
	uriWriteCloser.URI().Path()获取到的路径在安卓端与实际路径可能存在不符，	但是直接使用fyne.URIWriteCloser.Write方法可以写入。
	Windows端没有上述任何问题。
*/
func backupBthCallBack() {
	if true {
		dialog.ShowFileSave(func(uriWriteCloser fyne.URIWriteCloser, err error) {
			if uriWriteCloser == nil {
				// 保存对话框选择了取消
				return
			}
			r, jsonD, err := common.ReadFileAsBytes(storagejson.StoDPath)
			if !r {
				dialog.ShowInformation("err", "backupBthCallBack:"+err.Error(), storagejson.AppRef.W)
				return
			}
			_, err = uriWriteCloser.Write(jsonD)
			if err != nil {
				dialog.ShowInformation("err", "backupBthCallBack:"+err.Error(), storagejson.AppRef.W)
				return
			}
			dialog.ShowInformation("提示", "已备份数据到选择的目录中，可将其用于还原。", storagejson.AppRef.W)
			uriWriteCloser.Close()
		}, storagejson.AppRef.W)
		return
	}
	//dialog.ShowFileSave(func(uriWriteCloser fyne.URIWriteCloser, err error) {
	//	if uriWriteCloser == nil {
	//		// 保存对话框选择了取消
	//		return
	//	}
	//	var savePath string = uriWriteCloser.URI().Path()
	//	// 保存格式为json
	//	if !strings.HasSuffix(uriWriteCloser.URI().Name(), ".json") {
	//		savePath = uriWriteCloser.URI().Path() + ".json"
	//	}
	//	dialog.ShowInformation("1", savePath, storagejson.AppRef.W)
	//
	//	r, jsonD, err := common.ReadFileAsBytes(storagejson.StoDPath)
	//	if !r {
	//		dialog.ShowInformation("err", "backupBthCallBack:"+err.Error(), storagejson.AppRef.W)
	//		return
	//	}
	//
	//	r, err = common.CreateFile(savePath, jsonD)
	//	if !r {
	//		dialog.ShowInformation("err", "backupBthCallBack, common.CreateFile:"+err.Error(), storagejson.AppRef.W)
	//		return
	//	}
	//	dialog.ShowInformation("提示", "已备份数据到选择的目录中，可将其用于还原。", storagejson.AppRef.W)
	//	uriWriteCloser.Close()
	//}, storagejson.AppRef.W)

}

// 还原数据响应函数
/*
	安卓端必须使用fyne.URIReadCloser.Read将数据读取到缓存区，或许是fyne的问题，
	安卓端无法读取非app本地存储库中的文件，common.ReadFileAsBytes函数读取uriWriteCloser.URI().Path()提示目录不存在。
	uriWriteCloser.URI().Path()获取到的路径在安卓端与实际路径可能存在不符，
	但是直接使用fyne.URIReadCloser.Read方法可以读取非本地存储库的文件数据。
	Windows端没有上述任何问题。
*/
func restoreBthCallBack() {
	fileOpen := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if reader == nil {
			// 文件选择对话框选择了取消
			return
		}
		var jsonD []byte
		var buffer = make([]byte, 2048)
		// 从reader获取选择文件路径读取数据
		for true {
			i, err := reader.Read(buffer)
			if err != nil {
				// 如果读取到文件末尾，跳出循环
				if err.Error() == "EOF" {
					break
				}
				dialog.ShowInformation("err", "restoreBthCallBack, reader.Read:"+err.Error(), storagejson.AppRef.W)
				return
			}
			// 将本轮读取到的数据逐个追加到jsonD
			jsonD = append(jsonD, buffer[:i]...)
		}

		reader.Close()
		// 将还原的数据重新覆盖到本地存储库
		r, err := common.WriteExistedFile(storagejson.StoDPath, jsonD)
		if !r {
			dialog.ShowInformation("err", "restoreBthCallBack, common.WriteExistedFile:"+err.Error(), storagejson.AppRef.W)
			return
		}
		dialog.ShowInformation("提示", "数据已还原，请重启应用。", storagejson.AppRef.W)
		go func() {
			time.Sleep(time.Millisecond * 3000)
			storagejson.AppRef.W.Close()
		}()
	}, storagejson.AppRef.W)
	//fileOpen := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
	//	if reader == nil {
	//		// 文件选择对话框选择了取消
	//		return
	//	}
	//	// 从reader获取选择文件路径读取数据
	//	r, jsonD, err := common.ReadFileAsBytes(reader.URI().Path())
	//	if !r {
	//		dialog.ShowInformation("err", "restoreBthCallBack, common.ReadFileAsBytes:"+err.Error(), storagejson.AppRef.W)
	//		return
	//	}
	//	reader.Close()
	//	// 将还原的数据重新覆盖到本地存储库
	//	r, err = common.WriteExistedFile(storagejson.StoDPath, jsonD)
	//	if !r {
	//		dialog.ShowInformation("err", "restoreBthCallBack, common.WriteExistedFile:"+err.Error(), storagejson.AppRef.W)
	//		return
	//	}
	//	dialog.ShowInformation("提示", "数据已还原，请重启应用。", storagejson.AppRef.W)
	//	go func() {
	//		time.Sleep(time.Millisecond * 3000)
	//		storagejson.AppRef.W.Close()
	//	}()
	//}, storagejson.AppRef.W)
	//fileOpen.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	dialog.ShowConfirm("提示", "你即将还原备份数据，\r\n还原成功原本地数据将不可恢复，\r\n确定？", func(b bool) {
		if !b {
			return
		}
		fileOpen.Show()
	}, storagejson.AppRef.W)

}
