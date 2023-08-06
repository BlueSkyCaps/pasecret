package ui

import (
	"encoding/json"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"os"
	"pasecret/core/common"
	"pasecret/core/pi18n"
	"pasecret/core/preferences"
	"pasecret/core/storagedata"
	"strings"
	"time"
)

var treeSettingDictParent []string
var treeSettingDict map[string][]string

func createSettingTabContent() *widget.Tree {
	// 在函数中给treeSettingDictParent赋值，而不是全局。避免导包顺序LocalizedText()中localize变量还没被初始化
	treeSettingDictParent = []string{
		pi18n.LocalizedText("treeSettingDictParent-StartPwd", nil),
		pi18n.LocalizedText("treeSettingDictParent-DumpRestore", nil),
		pi18n.LocalizedText("treeSettingDictParent-Language", nil),
		pi18n.LocalizedText("treeSettingDictParent-Donate", nil),
		pi18n.LocalizedText("treeSettingDictParent-About", nil),
	}
	// 在函数中给treeSettingDict赋值，而不是全局。避免导包顺序LocalizedText()中localize变量还没被初始化
	treeSettingDict = map[string][]string{
		treeSettingDictParent[0]: {"NotChild_0"},
		treeSettingDictParent[1]: {
			pi18n.LocalizedText("treeSettingDict-Dump", nil),
			pi18n.LocalizedText("treeSettingDict-Restore", nil),
		},
		treeSettingDictParent[2]: {"NotChild_1"},
		treeSettingDictParent[3]: {"NotChild_3"},
		treeSettingDictParent[4]: {"NotChild_4"},
	}

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
				return
			}
			// 若是还原数据按钮
			if id == treeSettingDict[treeSettingDictParent[1]][1] {
				button.OnTapped = restoreBthCallBack
				return
			}
			// 若是启动密码
			if id == treeSettingDictParent[0] {
				button.OnTapped = lockPwdBthCallBack
				return
			}
			// 若是语言/language
			if id == treeSettingDictParent[2] {
				button.OnTapped = languageChiBthCallBack
				return
			}
			// 若是捐助赞赏
			if id == treeSettingDictParent[3] {
				button.OnTapped = donateBthCallBack
				return
			}
			// 若是关于
			if id == treeSettingDictParent[4] {
				button.OnTapped = aboutBthCallBack
				return
			}
			// 若是备份还原父按钮
			if id == treeSettingDictParent[1] {
				return
			}

			button.OnTapped = func() {
				dialog.ShowError(errors.New(pi18n.LocalizedText("settingTabShowErrorMaintenance", nil)),
					storagedata.AppRef.W)
			}
		})
	return tree
}

// 语言切换按钮回调
func languageChiBthCallBack() {
	lang := ""
	window := storagedata.AppRef.A.NewWindow("Language/语言")
	if !fyne.CurrentDevice().IsMobile() {
		window.Resize(fyne.Size{Height: 185, Width: 400})
	}
	vbox := container.NewVBox()
	radio := widget.NewRadioGroup([]string{"English", "中文"}, func(value string) {
		if value == "English" {
			lang = "en"
		} else if value == "中文" {
			lang = "zh"
		} else {
			lang = ""
		}
	})
	vbox.Add(radio)
	vbox.Add(widget.NewButton("Close/关闭", func() {
		window.Close()
	}))
	vbox.Add(widget.NewButton("OK/确定", func() {
		if common.IsWhiteAndSpace(lang) {
			return
		}
		// 更新首选项的语言值
		preferences.SetPreference("LocalLang", lang)
		dialog.ShowInformation(pi18n.LocalizedText("dialogShowInformationTitle", nil),
			"语言切换成功，请重启。\nswitch successful, please restart.", window)
		go func() {
			time.Sleep(time.Second * 3)
			os.Exit(0)
		}()
	}))
	window.SetContent(vbox)
	window.CenterOnScreen()
	window.Show()
}

// 启动密码 按钮回调
func lockPwdBthCallBack() {
	// 之前是否配置过启动密码
	ced := !common.IsWhiteAndSpace(preferences.GetPreferenceByLockPwd())
	window := storagedata.AppRef.A.NewWindow(pi18n.LocalizedText("setLockWindowTitle", nil))
	if !fyne.CurrentDevice().IsMobile() {
		window.CenterOnScreen()
		window.Resize(fyne.Size{Width: 300, Height: 400})
	}
	if ced {
		window.Show()
		dialog.ShowConfirm(pi18n.LocalizedText("lockPwdShowConfirmChi", nil),
			pi18n.LocalizedText("lockPwdAlreadySetShowConfirm", nil), func(b bool) {
				if b {
					preferences.RemovePreference("LockPwd")
					dialog.ShowInformation(pi18n.LocalizedText("dialogShowInformationTitle", nil),
						pi18n.LocalizedText("lockPwdClosedShowInformation", nil), window)
					go func() {
						time.Sleep(time.Millisecond * 3000)
						os.Exit(0)
					}()
				} else {
					window.Close()
				}
			}, window)
		return
	}
	tipLabel := widget.NewLabel(pi18n.LocalizedText("lockPwd4NumLabel", nil))
	pwdEntry := widget.NewPasswordEntry()
	common.EntryOnChangedEventHandler(pwdEntry)
	tipLabel2 := widget.NewLabel(pi18n.LocalizedText("lockPwdAgainConfirmLabel", nil))
	pwdEntry2 := widget.NewPasswordEntry()
	common.EntryOnChangedEventHandler(pwdEntry2)
	cBtn := widget.NewButton(pi18n.LocalizedText("lockPwdCancelLabel", nil), func() {
		window.Close()
	})
	yesBtn := widget.NewButton(pi18n.LocalizedText("lockPwdOkLabel", nil), func() {
		if !common.MatchPwdFormat(pwdEntry.Text) || !common.MatchPwdFormat(pwdEntry2.Text) {
			dialog.ShowInformation(pi18n.LocalizedText("dialogShowInformationTitle", nil),
				pi18n.LocalizedText("lockPwdNot4NumberShowInformation", nil), window)
			return
		}
		if pwdEntry.Text != pwdEntry2.Text {
			dialog.ShowInformation(pi18n.LocalizedText("dialogShowInformationTitle", nil),
				pi18n.LocalizedText("lockPwdNotNotMatchShowInformation", nil), window)
			return
		}
		dialog.ShowConfirm(pi18n.LocalizedText("dialogShowInformationTitle", nil),
			pi18n.LocalizedText("lockPwdSetTipShowConfirm", nil), func(b bool) {
				if b {
					preferences.SetPreference("LockPwd", pwdEntry.Text)
					dialog.ShowInformation(pi18n.LocalizedText("dialogShowInformationTitle", nil),
						pi18n.LocalizedText("lockPwdSetTipDoneInformation", nil), window)
					go func() {
						time.Sleep(time.Millisecond * 3000)
						os.Exit(0)
					}()
				}
			}, window)
	})

	center := container.NewCenter(container.NewVBox(tipLabel, pwdEntry, tipLabel2, pwdEntry2, cBtn, yesBtn))
	window.SetContent(center)
	window.Show()
}

// 关于 按钮回调响应
func aboutBthCallBack() {
	window := storagedata.AppRef.A.NewWindow(pi18n.LocalizedText("aboutWindowTitle", nil))
	appLinkUri, err := url.Parse(common.AppLinkUri_)
	githubUri, err := url.Parse(common.GithubUri)
	blogUri, err := url.Parse(common.BlogUri)
	if err != nil {
		return
	}
	statementLabel := widget.NewLabel("")
	statementLabel.SetText(pi18n.LocalizedText("aboutStatementLabelText", nil))
	statementScroll := container.NewHScroll(statementLabel)
	statementScroll.Hide()
	center := container.NewCenter(container.NewVBox(
		container.NewHScroll(widget.NewLabelWithStyle(pi18n.LocalizedText("aboutIntroduceLabelText", nil),
			fyne.TextAlignCenter, fyne.TextStyle{Bold: true})),
		container.NewCenter(container.NewHBox(
			widget.NewHyperlink(pi18n.LocalizedText("aboutAppSiteLinkName", nil), appLinkUri),
			widget.NewLabel("-"),
			widget.NewHyperlink("Github", githubUri),
			widget.NewLabel("-"),
			widget.NewHyperlink(pi18n.LocalizedText("aboutMySiteLinkName", nil), blogUri),
		)),
		widget.NewButton(pi18n.LocalizedText("aboutStatementButtonText", nil), func() {
			statementScroll.Show()
		}),
		statementScroll,
	))
	window.SetContent(center)
	window.Show()
}

// 捐助赞赏 按钮响应函数
func donateBthCallBack() {
	window := storagedata.AppRef.A.NewWindow(pi18n.LocalizedText("donateWindowTitle", nil))
	if !fyne.CurrentDevice().IsMobile() {
		window.Resize(fyne.Size{Height: 700})
	}
	uriWechat, err := storage.ParseURI(common.DonateWechatUri_)
	uriAlipay, err := storage.ParseURI(common.DonateAliPayUri_)
	if err != nil {
		dialog.ShowInformation("err", "donateBthCallBack:"+err.Error(), storagedata.AppRef.W)
		return
	}
	wechatimage := canvas.NewImageFromURI(uriWechat)
	alipayimage := canvas.NewImageFromURI(uriAlipay)
	// 窗口将会自动达到图片原图尺寸大小，但有时候会显示不全的bug，设置window.Resize
	wechatimage.FillMode = canvas.ImageFillOriginal
	alipayimage.FillMode = canvas.ImageFillOriginal
	box := container.NewVBox(widget.NewLabel("微信扫一扫"), wechatimage, widget.NewLabel("支付宝扫一扫"), alipayimage)
	if fyne.CurrentDevice().IsMobile() {
		// 移动端因为尺寸小，可能达不到图片显示的高度，会被截停，因此添加垂直滚动条
		window.SetContent(container.NewVScroll(box))
	} else {
		window.SetContent(box)
	}
	dialog.ShowConfirm(pi18n.LocalizedText("dialogShowInformationTitle", nil),
		pi18n.LocalizedText("donateOpenShowConfirm", nil), func(b bool) {
			if !b {
				return
			}
			window.Show()
		}, storagedata.AppRef.W)
}

// 备份数据 响应函数
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
			r, jsonD, err := common.ReadFileAsBytes(storagedata.StoDPath)
			if !r {
				dialog.ShowInformation("err", "backupBthCallBack:"+err.Error(), storagedata.AppRef.W)
				return
			}
			_, err = uriWriteCloser.Write(jsonD)
			if err != nil {
				dialog.ShowInformation("err", "backupBthCallBack:"+err.Error(), storagedata.AppRef.W)
				return
			}
			dialog.ShowInformation(pi18n.LocalizedText("dialogShowInformationTitle", nil),
				pi18n.LocalizedText("dumpDoneShowInformation", nil), storagedata.AppRef.W)
			uriWriteCloser.Close()
		}, storagedata.AppRef.W)
		return
	}

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
				dialog.ShowInformation("err", "restoreBthCallBack, reader.Read:"+err.Error(), storagedata.AppRef.W)
				return
			}
			// 将本轮读取到的数据逐个追加到jsonD
			jsonD = append(jsonD, buffer[:i]...)
		}
		reader.Close()
		err = json.Unmarshal(jsonD, &storagedata.AppRef.LoadedItems)
		if err != nil {
			dialog.ShowCustom(pi18n.LocalizedText("dialogShowInformationTitle", nil),
				pi18n.LocalizedText("restoreInvalidShowCustom", map[string]interface{}{"errMsg": err.Error()}),
				widget.NewLabel(""), storagedata.AppRef.W)
			return
		}
		// 将还原的数据重新覆盖到本地存储库
		r, err := common.WriteExistedFile(storagedata.StoDPath, jsonD)
		if !r {
			dialog.ShowInformation("err", "restoreBthCallBack, common.WriteExistedFile:"+err.Error(),
				storagedata.AppRef.W)
			return
		}
		dialog.ShowInformation(pi18n.LocalizedText("dialogShowInformationTitle", nil),
			pi18n.LocalizedText("restoreDoneShowInformation", nil), storagedata.AppRef.W)
		go func() {
			time.Sleep(time.Millisecond * 3000)
			//移动端无法退出：storagedata.AppRef.A.Quit()
			os.Exit(0)
		}()

	}, storagedata.AppRef.W)
	dialog.ShowConfirm(pi18n.LocalizedText("dialogShowInformationTitle", nil),
		pi18n.LocalizedText("restoreBeginShowConfirm", nil),
		func(b bool) {
			if !b {
				return
			}
			fileOpen.Show()
		}, storagedata.AppRef.W)

}
