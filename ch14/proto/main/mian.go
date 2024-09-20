package main

import (
	"encoding/json"
	"fmt"
	"go-learn/ch14/proto/helloworld"

	"google.golang.org/protobuf/proto"
)

type Hello struct {
	Name    string   `json:"name`
	Age     int      `json:"age"`
	Courses []string `json:"courses"`
}

func main() {
	req := helloworld.HelloRequest{
		Name:    "ali",
		Age:     18,
		Courses: []string{"go", "python", "java"},
	}

	rsp, _ := proto.Marshal(&req)

	jsonStruct := Hello{
		Name:    "ali",
		Age:     18,
		Courses: []string{"go", "python", "java"},
	}
	jsonRes, _ := json.Marshal(jsonStruct)
	fmt.Println(string(rsp))
	fmt.Println(string(jsonRes))
}
