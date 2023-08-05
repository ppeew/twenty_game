package tests

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func TestReadImage(t *testing.T) {
	filePath := "C:\\Users\\22378\\GolandProjects\\twenty_game\\srv\\user_srv\\images\\123.jpg"
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	var ret []byte
	for true {
		p := make([]byte, 1024)
		n, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				//读到结尾
				break
			} else {
				panic(err)
			}
		}
		ret = append(ret, p[:n]...)
	}
	open, _ := os.OpenFile("C:\\Users\\22378\\GolandProjects\\twenty_game\\srv\\user_srv\\images\\make.jpg", os.O_CREATE, 0666)
	writer := bufio.NewWriter(open)
	writer.Write(ret)
	open.Close()
}
