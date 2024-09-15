package main

import (
	"fmt"
	. "go-learn/ch10/user"
	course "go-learn/ch10/user"

	// _ 匿名导入 包里如有 init 函数会自动执行
	_ "go-learn/ch09/user"
)

func main() {
	c := course.Course{
		Name: "go",
	}

	fmt.Println(course.GetCourse(c))

	c1 := Course{
		Name: "go1",
	}

	fmt.Println(GetCourse(c1))
}
