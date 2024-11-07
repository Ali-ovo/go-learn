package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func (s *Student) SetName(name string) error {
	if name == "" || s.Name == "" {
		return errors.New("name can't be empty")
	}
	s.Name = name
	return nil
}

type Student struct {
	Name string
	Age  int
}

func NewStudent() (*Student, error) {
	stu := &Student{
		Age: 18,
	}

	e := stu.SetName("")
	if e != nil {
		return stu, errors.Wrap(e, "set name")
	}
	return stu, nil
}

func main() {
	_, e := NewStudent()
	fmt.Printf("%+v\n", e)
}
