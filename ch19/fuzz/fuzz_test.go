package fuzz

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func Reverse(s string) (string, error) {
	//没有考虑到非法的unicode编码
	//q 代表该值对应的单引号括起来的go语法字符字面值
	if !utf8.ValidString(s) {
		return "", fmt.Errorf("invalid utf8: %q", s)
	}
	b := []rune(s)
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b), nil
}

func TestReverse(t *testing.T) {
	testcases := []struct {
		in, want string
	}{
		{"hello", "olleh"}, //基本
		{"a", "a"},         //边界
		{" ", " "},         //特殊
	}
	for _, c := range testcases {
		rev, _ := Reverse(c.in)
		assert.Equal(t, c.want, rev)
	}
}
