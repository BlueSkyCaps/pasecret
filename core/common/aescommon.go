package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var appDefaultKeyAES = "theKeyIsJustHereNotUsedOnProduct"

// KeyBytesAES 生成密钥（256位）的字节形式
func KeyBytesAES(k string) ([]byte, error) {
	key := []byte(k)
	if IsWhiteAndSpace(k) {
		a := appDefaultKeyAES
		key = []byte(a)
		return key, nil
	}
	if len([]byte(k)) != 32 {
		return nil, errors.New("bytes size must 32")
	}
	return key, nil
}

// EncryptAES 用AES密钥加密明文生成密文
func EncryptAES(key []byte, plaintext string) (string, error) {
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

// DecryptAES 用AES密钥解密密文返回明文
func DecryptAES(key []byte, ciphertext string) (string, error) {
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
