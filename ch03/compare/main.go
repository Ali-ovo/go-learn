package main

import (
	"fmt"
	"strings"
)

func main() {
	// 字符串比较
	str1 := "hello"
	str2 := "bello"
	if str1 == str2 {
		println("str1 == str2")
	}

	fmt.Println(str1 > str2)

	strLen := len(str1)
	fmt.Println(strLen)

	name := " ali study go,lalala "
	fmt.Printf("strings.Contains:  %v\n", strings.Contains(name, "go"))

	// 字符串出现次数
	fmt.Printf("strings.Count: %v\n", strings.Count(name, "a"))

	// 分割
	fmt.Printf("strings.Split: %v\n", strings.Split(name, ","))

	// 是否包含前缀后缀
	fmt.Printf("strings.HasPrefix: %v\n", strings.HasPrefix(name, "ali"))

	// 查找子串位置
	fmt.Printf("strings.Index: %v\n", strings.Index(name, "go"))
	fmt.Printf("strings.Rune: %v\n", strings.IndexRune(name, []rune(name)[10]))

	// 子串替换  -1 代表全部替换
	fmt.Printf("strings.Replace: %v\n", strings.Replace(name, "la", "hahala", -1))
	fmt.Printf("strings.ReplaceAll: %v\n", strings.ReplaceAll(name, "la", "jjla"))

	// 大小写转换
	fmt.Printf("strings.ToLower: %v\n", strings.ToLower("GO"))
	fmt.Printf("strings.ToUpper: %v\n", strings.ToUpper("go"))

	// 去掉空格
	fmt.Printf("strings.TrimSpace: %v\n", strings.TrimSpace(name))

	// 去掉两侧的特殊字符
	fmt.Printf("strings.Trim: %v\n", strings.Trim(name, " a"))
}
