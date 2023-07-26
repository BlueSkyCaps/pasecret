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
	plaintext := "Hello, AES2222wd是，，，发唧唧复唧唧大约为七点七一缗彂!"

	// 加密数据
	encrypted, err := encrypt(key, plaintext)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}

	// 解密数据
	decrypted, err := decrypt(key, encrypted)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}

	fmt.Println("Original:", plaintext)
	fmt.Println("Encrypted:", encrypted)
	fmt.Println("Decrypted:", decrypted)
}
