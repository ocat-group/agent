package util

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// GetMD5 获取byte对应MD5
func GetMD5(bytes []byte) string {
	m := md5.New()
	m.Write(bytes)
	return hex.EncodeToString(m.Sum(nil))
}

// ReadFileMd5 获取文件byte
func ReadFileMd5(file string) (string, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return GetMD5(bytes), nil
}
