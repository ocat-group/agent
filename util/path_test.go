package util

import (
	"log"
	"testing"
)

func TestAppendPathSeparator(t *testing.T) {
	path := AppendPathSeparator("D:\\DevSoftware\\Scripte")
	log.Println(path)
} // grpc支持查询程序信息接口
