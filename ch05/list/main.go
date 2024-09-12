package main

import (
	"container/list"
	"fmt"
)

func main() {
	// list
	myList := list.New()
	myList.PushBack(1)
	myList.PushBack(2)

	myList.PushFront(0)

	// 正序遍历
	i := myList.Front()
	for ; i != nil; i = i.Next() {
		if i.Value == 2 {
			myList.InsertBefore(1.5, i)
		}
		fmt.Println(i.Value)
	}

	// 倒序遍历
	for i := myList.Back(); i != nil; i = i.Prev() {
		fmt.Println(i.Value)
	}
}
