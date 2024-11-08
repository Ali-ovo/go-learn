package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type DBError struct {
	msg string
}

func (d *DBError) Error() string {
	return d.msg
}

var ErrNameEmpty = &DBError{"Name is empty"}

func (s *Student) SetName(name string) error {
	if name == "" || s.Name == "" {
		return ErrNameEmpty
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
	var perr *DBError
	if errors.As(e, *perr) {
		fmt.Println("Name is empty")
	}
}
