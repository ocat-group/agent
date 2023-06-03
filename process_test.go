package main

import (
	"log"
	"testing"
)

func TestAppendPathSeparator(t *testing.T) {
	path := appendPathSeparator("D:\\DevSoftware\\Scripte")
	log.Println(path)
}
