package main

import (
	"agent/plugin_manager"
	"log"
	"testing"
)

func TestAppendPathSeparator(t *testing.T) {
	path := plugin_manager.AppendPathSeparator("D:\\DevSoftware\\Scripte")
	log.Println(path)
} // grpc支持查询程序信息接口
