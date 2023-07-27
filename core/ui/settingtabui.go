package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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
func backupBthCallBack() {
	dialog.ShowFileSave(func(uriWriteCloser fyne.URIWriteCloser, err error) {
		var savePath string = uriWriteCloser.URI().Path()
		if uriWriteCloser == nil {
			// 保存对话框选择了取消
			return
		}
		// 保存格式为json
		if !strings.HasSuffix(uriWriteCloser.URI().Name(), ".json") {
			savePath = uriWriteCloser.URI().Path() + ".json"
		}

		r, jsonD, err := common.ReadFileAsBytes(storagejson.StoDPath)
		if !r {
			dialog.ShowInformation("err", "backupBthCallBack:"+err.Error(), storagejson.AppRef.W)
			return
		}
		r, err = common.CreateFile(savePath, jsonD)
		if !r {
			dialog.ShowInformation("err", "backupBthCallBack, common.CreateFile:"+err.Error(), storagejson.AppRef.W)
			return
		}
		dialog.ShowInformation("提示", "已备份数据到选择的目录中，可将其用于还原。\r\n"+
			"文件完整路径为："+savePath, storagejson.AppRef.W)
		uriWriteCloser.Close()
	}, storagejson.AppRef.W)

}

// 还原数据响应函数
func restoreBthCallBack() {
	fileOpen := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if reader == nil {
			// 文件选择对话框选择了取消
			return
		}
		reader.Close()
		// 从reader获取选择文件路径读取数据
		r, jsonD, err := common.ReadFileAsBytes(reader.URI().Path())
		if !r {
			dialog.ShowInformation("err", "restoreBthCallBack, common.ReadFileAsBytes:"+err.Error(), storagejson.AppRef.W)
			return
		}
		// 将还原的数据重新覆盖到本地存储库
		r, err = common.WriteExistedFile(storagejson.StoDPath, jsonD)
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
	fileOpen.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	dialog.ShowConfirm("提示", "你即将还原备份数据，还原成功原本地数据将不可恢复，确定？", func(b bool) {
		if !b {
			return
		}
		fileOpen.Show()
	}, storagejson.AppRef.W)

}
