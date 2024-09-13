package main

import (
	"errors"
	"fmt"
)

func errorFn() (int, error) {

	return 0, errors.New("this is error")
}

func main() {
	// error

	_, err := errorFn()
	if err != nil {
		fmt.Println(err)
	}

}
