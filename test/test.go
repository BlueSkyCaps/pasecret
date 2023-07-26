package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"pasecret/core/common"
)

func encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(key []byte, ciphertext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(decoded) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := decoded[:aes.BlockSize]
	decoded = decoded[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(decoded, decoded)

	return string(decoded), nil
}

func main() {

	key, err := common.KeyBytesAES(common.AppProductKeyAES)
	if err != nil {
		println(err.Error())
	}
	// 待加密的数据

	name := "In+lnVpEiA/qQ8CBbTgpaIScIbPoBQeLKoleyQ=="
	// 加密数据
	name, err = decrypt(key, name)
	println(name)
	accountName := "cGqeKMBARnIeAI57bPEZ6CqmEUGyUIYJ"
	// 加密数据
	accountName, err = decrypt(key, accountName)
	println(accountName)
	password := "HN0hZAW1PXb0cbK7nw3dGNiz"
	// 加密数据
	password, err = decrypt(key, password)
	println(password)
	site := "JjxX6m0F2nC27P1Nqez36t1A2bgQOQ9MOLc/NW/K3wFaVasVNSTOQk4="
	// 加密数据
	site, err = decrypt(key, site)
	println(site)
	remark := "R7TjREKMX3s06daubeZZdMffSl2OGYRQaDqxMRmieuNFyICBJ8MegqOCNv3Xi97gF0ZzESE85ksGVnBoAveXypd3YKsiTW6T+0pivtinHyG6n7fbMMrLF1Y3SYU33cDBK7noEaXt/lTv+H5gvw=="
	// 加密数据
	remark, err = decrypt(key, remark)
	println(remark)

}
