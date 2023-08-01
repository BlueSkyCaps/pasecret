// Package common 共用函数
package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Existed 目录或文件是否存在
func Existed(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CreateDir 创建路径中的目录
func CreateDir(path string) (bool, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err == nil {
		return true, nil
	}
	return false, err
}

// CreateFile 创建路径文件，并可写入文件数据
func CreateFile(path string, data []byte) (bool, error) {
	// Android必须添加os.O_WRONLY，否则报"bad file descriptor"权限问题
	create, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return false, err
	}
	_, err = create.Write(data)
	if err != nil {
		create.Close()
		return false, err
	}
	err = create.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

// WriteExistedFile 已有路径文件覆盖写入数据
func WriteExistedFile(path string, data []byte) (bool, error) {
	// Android必须添加os.O_WRONLY，否则报"bad file descriptor"权限问题
	create, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return false, err
	}
	_, err = create.Write(data)
	if err != nil {
		create.Close()
		return false, err
	}
	err = create.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

// ReadFileAsString 从路径文件读全部数据，字符串形式
func ReadFileAsString(path string) (bool, string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false, "", err
	}
	return true, string(b), nil
}

// ReadFileAsBytes 从路径文件读全部数据，byte形式
func ReadFileAsBytes(path string) (bool, []byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false, nil, err
	}
	return true, b, nil

}

type EditForm struct {
	Name        string
	Alias       string
	Description string
}

func IsWhiteAndSpace(s string) bool {
	return strings.TrimSpace(s) == ""
}

// GenAscRankId 根据当前日期生成一个整型数字，用于归类夹升序排序。每次生成的值都比之前生成的要大，最小值基于秒
func GenAscRankId() (bool, int) {
	now := time.Now()
	//调用结构体中的方法：
	rankS := fmt.Sprintf("%v%v%v%v%v%v", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second())
	rank, err := strconv.Atoi(rankS)
	if err != nil {
		return false, 0
	}
	return true, rank
}

// SearchDataResultViewModel 用于关键字搜索密码项展示table的视图模型
type SearchDataResultViewModel struct {
	VDataName        string
	VDataAccountName string
	VDataPassword    string
	VDataSite        string
	VDataRemark      string
	//VDataId           int
	VDataCategoryName string
}

// SearchDataResultHeader 用于关键字搜索密码项展示table位于表头的列明
var SearchDataResultHeader = SearchDataResultViewModel{
	VDataName:        "项名称",
	VDataAccountName: "账号",
	VDataPassword:    "密码",
	VDataSite:        "网址",
	VDataRemark:      "备注",
	//VDataId           0
	VDataCategoryName: "所属归类夹",
}

// ShowSearchResultCeilWH_ 用于设置关键字搜索密码项展示table的单元格长度。
//单元格宽高度可以由此填充字符定义，不要使用canvas.Text无法文本换行，不建议使用SetColumnWidth等设置宽高
var ShowSearchResultCeilWH_ = "xxxxxxxxxxxxxxxxxxxxxxxx\nxxxxxxxxxxxxxxxxxxxxxxxx\nxxxxxxxxxxxxxxxxxxxxxxxx"
var DonateAliPayUri_ = "http://sto.reminisce.top/OTHER/donate/alipay.jpg"
var DonateWechatUri_ = "http://sto.reminisce.top/OTHER/donate/wechat.jpg"
var AppLinkUri_ = "https://apps.reminisce.top/pasecret"
var GithubUri = "https://github.com/BlueSkyCaps/pasecret"
var BlogUri = "https://www.reminisce.top/"

func MatchPwdFormat(str string) bool {
	// 定义一个表示4个数字的正则表达式
	re, _ := regexp.Compile(`^\d{4}$`)

	// 使用正则表达式进行匹配
	if re.MatchString(str) {
		return true
	}
	return false
}
