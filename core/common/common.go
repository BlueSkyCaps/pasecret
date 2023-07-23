// Package common 共用函数
package common

import (
	"io/ioutil"
	"os"
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
