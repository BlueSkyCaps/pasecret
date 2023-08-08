package common

import (
	"crypto/sha256"
	"encoding/hex"
)

// EncryptSha 生成sha-256散列值
func EncryptSha(plaintext string) string {
	hash := sha256.Sum256([]byte(plaintext))
	sha256String := hex.EncodeToString(hash[:])
	return sha256String
}
