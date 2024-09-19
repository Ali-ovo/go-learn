package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kirinlabs/HttpRequest"
)

type ResponseData struct {
	Data int `json:"data"`
}

func Add(a, b int) int {
	req := HttpRequest.NewRequest()
	req.SetTimeout(20 * time.Second)

	url := fmt.Sprintf("http://127.0.0.1:8000/add?a=%d&b=%d", a, b)
	res, _ := req.Get(url, nil)

	body, _ := res.Body()
	fmt.Println(string(body))
	respData := ResponseData{}

	_ = json.Unmarshal(body, &respData)

	return respData.Data
}

func main() {
	fmt.Println(Add(2, 2))
}
