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
	"pasecret/core/preferences"
	"pasecret/core/storagejson"
	"strings"
	"time"
)

var treeSettingDictParent = []string{"启动密码", "备份还原", "多端同步", "捐助赞赏", "关于"}
var treeSettingDict map[string][]string = map[string][]string{
	treeSettingDictParent[0]: {"NotChild_0"},
	treeSettingDictParent[1]: {"备份数据", "还原数据到本机"},
	treeSettingDictParent[2]: {"NotChild_1"},
	treeSettingDictParent[3]: {"NotChild_3"},
	treeSettingDictParent[4]: {"NotChild_4"},
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
				dialog.ShowError(errors.New("维护中，尽情期待。"), storagejson.AppRef.W)
			}
		})
	return tree
}

// 启动密码 按钮回调
func lockPwdBthCallBack() {
	// 之前是否配置过启动密码
	ced := !common.IsWhiteAndSpace(preferences.GetPreferenceByLockPwd())
	window := storagejson.AppRef.A.NewWindow("设置启动密码")
	if !fyne.CurrentDevice().IsMobile() {
		window.CenterOnScreen()
		window.Resize(fyne.Size{Width: 300, Height: 400})
	}
	if ced {
		window.Show()
		dialog.ShowConfirm("选择", "您已经设置过启动密码。\n您需要关闭或者重设密码吗？", func(b bool) {
			if b {
				preferences.RemovePreferenceByLockPwd()
				dialog.ShowInformation("提示", "已关闭启动密码。\n您现在可以选择重新设置了。", window)
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
	tipLabel := widget.NewLabel("输入四位数密码：")
	pwdEntry := widget.NewPasswordEntry()
	common.EntryOnChangedEventHandler(pwdEntry)
	tipLabel2 := widget.NewLabel("再次确认输入：")
	pwdEntry2 := widget.NewPasswordEntry()
	common.EntryOnChangedEventHandler(pwdEntry2)
	cBtn := widget.NewButton("取消", func() {
		window.Close()
	})
	yesBtn := widget.NewButton("确定", func() {
		if !common.MatchPwdFormat(pwdEntry.Text) || !common.MatchPwdFormat(pwdEntry2.Text) {
			dialog.ShowInformation("提示", "请输入4个数字！", window)
			return
		}
		if pwdEntry.Text != pwdEntry2.Text {
			dialog.ShowInformation("提示", "两次密码不匹配！", window)
			return
		}
		dialog.ShowConfirm("重要提示！",
			"准备设置启动密码:\n"+
				"请您牢记您刚刚设置的4位数字，\n"+
				"若您忘记，将无法进入应用。\n"+
				"建议您先备份数据，若您遗忘密码，\n"+
				"可以尽可能保留存储的数据。\n"+
				"最后，是否确定开启？！", func(b bool) {
				if b {
					preferences.SetPreferenceByLockPwd(pwdEntry.Text)
					dialog.ShowInformation("提示", "已开启，请重启应用。", window)
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
	window := storagejson.AppRef.A.NewWindow("关于")
	appLinkUri, err := url.Parse(common.AppLinkUri_)
	githubUri, err := url.Parse(common.GithubUri)
	blogUri, err := url.Parse(common.BlogUri)
	if err != nil {
		return
	}
	statementLabel := widget.NewLabel("")
	statementLabel.SetText(
		"本软件作者：BlueSkyCaps。本软件不传输任何数据，只有在“捐助赞赏”中需要联网加载付款二维码，\n" +
			"并且事先有提示是否打开。本软件加密存储您保存的数据，但不代表百分百能够保障您的数据安全。\n" +
			"如您在使用此软件过程中产生数据丢失、账户密码泄露造成的损失，本软件和作者不负任何责任，\n" +
			"损失由您自己或其他方面造成且承担。\n" +
			"使用本软件代表您同意此内容。")
	statementScroll := container.NewHScroll(statementLabel)
	statementScroll.Hide()
	center := container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Pasecret是一款能在多个平台运行的账号密码管理软件。\n"+
			"例如，您可以在手机上使用，并且同步数据到电脑端。\n"+
			"数据采用加密算法，并且可以断网使用，不会进行远程传输。\n"+
			"欢迎进行捐助赞赏，作者表示感激。",
			fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(container.NewHBox(
			widget.NewHyperlink("软件主页", appLinkUri),
			widget.NewLabel("-"),
			widget.NewHyperlink("Github", githubUri),
			widget.NewLabel("-"),
			widget.NewHyperlink("作者博客", blogUri),
		)),
		widget.NewButton("免责声明等", func() {
			statementScroll.Show()
		}),
		statementScroll,
	))
	window.SetContent(center)
	window.Show()
}

// 捐助赞赏 按钮响应函数
func donateBthCallBack() {
	window := storagejson.AppRef.A.NewWindow("捐助赞赏")
	if !fyne.CurrentDevice().IsMobile() {
		window.Resize(fyne.Size{Height: 800})
	}
	uriWechat, err := storage.ParseURI(common.DonateWechatUri_)
	uriAlipay, err := storage.ParseURI(common.DonateAliPayUri_)
	if err != nil {
		dialog.ShowInformation("err", "donateBthCallBack:"+err.Error(), storagejson.AppRef.W)
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
	dialog.ShowConfirm("嘿嘿", "捐助赞赏需要进行网络连接\n"+
		"获取赞赏二维码、赞赏者列表。\n本应用不会传递任何其他数据，\n是否继续？", func(b bool) {
		if !b {
			return
		}
		window.Show()
	}, storagejson.AppRef.W)
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
		err = json.Unmarshal(jsonD, &storagejson.AppRef.LoadedItems)
		if err != nil {
			dialog.ShowCustom("错误", "还原失败，不是有效的Pasecret数据文件！\n"+err.Error(),
				widget.NewLabel(""), storagejson.AppRef.W)
			return
		}
		// 将还原的数据重新覆盖到本地存储库
		r, err := common.WriteExistedFile(storagejson.StoDPath, jsonD)
		if !r {
			dialog.ShowInformation("err", "restoreBthCallBack, common.WriteExistedFile:"+err.Error(),
				storagejson.AppRef.W)
			return
		}
		dialog.ShowInformation("提示", "数据已还原，请重新打开程序。", storagejson.AppRef.W)
		go func() {
			time.Sleep(time.Millisecond * 3000)
			//移动端无法退出：storagejson.AppRef.A.Quit()
			os.Exit(0)
		}()

	}, storagejson.AppRef.W)
	dialog.ShowConfirm("提示", "你即将还原备份数据，\n还原成功原本地数据将不可恢复，\r\n确定？", func(b bool) {
		if !b {
			return
		}
		fileOpen.Show()
	}, storagejson.AppRef.W)

}
