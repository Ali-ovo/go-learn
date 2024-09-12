package main

import "fmt"

func main() {
	// map
	map1 := map[string]string{
		"web": "react",
		"go":  "gin",
	}

	fmt.Println(map1)

	map1["web"] = "vue"
	fmt.Println(map1)

	// 必须初始化 不然无法放值
	map2 := map[string]string{}
	map3 := make(map[string]string, 1)
	fmt.Println(map2, map3)

	// 遍历  是无序的
	for key, value := range map1 {
		fmt.Println(key, value)
	}
	// 判断是否存在
	if _, ok := map1["web"]; ok {
		fmt.Println("key web is exist")
	}

	// 删除
	delete(map1, "web")
	fmt.Println(map1)
}
