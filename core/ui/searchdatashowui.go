package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"pasecret/core/common"
	"pasecret/core/pi18n"
	"pasecret/core/storagedata"
)

var SearchDataResultHeader common.SearchDataResultViewModel

// ShowSearchResultWin 点击了搜索按钮，根据文本框关键词显示结果
func ShowSearchResultWin(text string) {
	// 在函数中给SearchDataResultHeader赋值，而不是全局。避免导包顺序LocalizedText()中localize变量还没被初始化
	SearchDataResultHeader = common.SearchDataResultViewModel{
		VDataName:        pi18n.LocalizedText("SearchDataResultHeaderVDataName", nil),
		VDataAccountName: pi18n.LocalizedText("SearchDataResultHeaderVDataAccountName", nil),
		VDataPassword:    pi18n.LocalizedText("SearchDataResultHeaderVDataPassword", nil),
		VDataSite:        pi18n.LocalizedText("SearchDataResultHeaderVDataSite", nil),
		VDataRemark:      pi18n.LocalizedText("SearchDataResultHeaderVDataRemark", nil),
		//VDataId           0
		VDataCategoryName: pi18n.LocalizedText("SearchDataResultHeaderVDataCategoryName", nil),
	}
	// 根据关键词获取数据视图模型
	viewModels := storagedata.SearchByKeywordTo(text)
	resW := storagedata.AppRef.A.NewWindow(pi18n.LocalizedText("SearchWindowTitle", nil))
	// 将表头元素放在首位
	viewModels = append(viewModels, common.SearchDataResultViewModel{})
	copy(viewModels[1:], viewModels)
	viewModels[0] = SearchDataResultHeader
	// 创建显示搜索结果的表格
	table := widget.NewTable(
		func() (int, int) {
			return len(viewModels), 6
		},
		func() fyne.CanvasObject {
			// 单元格宽高度可以由此填充字符定义，不要使用canvas.Text无法文本换行，不建议使用SetColumnWidth等设置宽高
			text := widget.NewLabel(common.ShowSearchResultCeilWH_)
			return text
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			currentRowIndex := i.Row
			currentRowVm := viewModels[currentRowIndex]

			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(currentRowVm.VDataCategoryName)
			case 1:
				o.(*widget.Label).SetText(currentRowVm.VDataName)
			case 2:
				o.(*widget.Label).SetText(currentRowVm.VDataAccountName)
			case 3:
				o.(*widget.Label).SetText(currentRowVm.VDataPassword)
			case 4:
				o.(*widget.Label).SetText(currentRowVm.VDataSite)
			case 5:
				o.(*widget.Label).SetText(currentRowVm.VDataRemark)

			}
		})
	// 设置行高
	//for i := 0; i < len(viewModels); i++ {
	//	if i == 0 {
	//		table.SetRowHeight(i, 50)
	//	} else {
	//		table.SetRowHeight(i, 30)
	//	}
	//}
	//table.SetColumnWidth(0, 100)
	//table.SetColumnWidth(1, 100)
	//table.SetColumnWidth(2, 100)
	//table.SetColumnWidth(3, 100)
	//table.SetColumnWidth(4, 100)
	//table.SetColumnWidth(5, 100)
	resW.Resize(fyne.Size{Width: 630, Height: 500})
	resW.SetContent(table)
	resW.CenterOnScreen()
	resW.RequestFocus()
	resW.Show()
}
