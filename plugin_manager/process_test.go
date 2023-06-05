package plugin_manager

import (
	"log"
	"testing"
)

func TestAppendPathSeparator(t *testing.T) {
	path := appendPathSeparator("D:\\DevSoftware\\Scripte")
	log.Println(path)
} // grpc支持查询程序信息接口
