package main

import "fmt"

func main() {
	// slice
	courses := []string{"go", "python", "java"}
	courses = append(courses, "php", "c++")
	fmt.Println(courses)

	// make 创建切片
	allCourses := [3]string{"go", "python", "java"}
	coursesSlice := allCourses[0:]
	fmt.Println(coursesSlice)

	makeSlice := make([]string, 3)
	fmt.Println(makeSlice)

	// 批量添加
	coursesSlice2 := []string{"web", "c#"}
	coursesSlice = append(coursesSlice, coursesSlice2[:1]...)
	fmt.Println(coursesSlice)

	// 删除
	coursesSlice = append(coursesSlice[:1], coursesSlice[2:]...)
	fmt.Println(coursesSlice)

	// 复制
	copySlice := coursesSlice[:]
	fmt.Println(copySlice)

	var copySlice2 = make([]string, len(coursesSlice))

	// copy 会覆盖
	copy(copySlice2, coursesSlice)
	fmt.Println(copySlice2)
}
