package main

import (
	"errors"
	"fmt"
)

func errorFn() (int, error) {
	// recover 这个函数可以捕获 panic
	// recover 只能在 defer 中使用
	// recover 处理异常后逻辑代码并不会恢复到 panic 中
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err, 999)
		}
	}()

	// panic 会导致程序退出 平时开发中不要随便使用
	panic("this is an painc")

	return 0, errors.New("this is error")
}

func main() {
	// error

	_, err := errorFn()
	if err != nil {
		fmt.Println(err)
	}

}
