package tests

import (
	"strings"
	"testing"
	"time"
)

func TestFile(t *testing.T) {
	png := "1.png"
	filePathByte, _ := time.Now().MarshalText()
	filePath := string(filePathByte)
	filePath = strings.Split(filePath, ".")[0]
	filePath = strings.Replace(filePath, "T", "/", 1) + "-" + png
	println(filePath)
}
