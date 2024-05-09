package tests

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestDecode(t *testing.T) {
	str := "\xE4\xBB\x81" //中文 仁
	r, size := utf8.DecodeRuneInString(str)
	fmt.Println(string(r), "dd", size)
}
