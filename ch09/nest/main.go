package main

import "fmt"

type MyWriter interface {
	Write(string)
}

type MyReader interface {
	Read() string
}

type MyReadWriter interface {
	MyReader
	MyWriter
	ReadWrite()
}

type SreadWriter struct{}

// Read implements MyReadWriter.
func (s *SreadWriter) Read() string {
	fmt.Println("read")
	return "read"
}

// ReadWrite implements MyReadWriter.
func (s *SreadWriter) ReadWrite() {

	fmt.Println("read write")
}

// Write implements MyReadWriter.
func (s *SreadWriter) Write(string) {

	fmt.Println("write")
}

func main() {
	var rw MyReadWriter = &SreadWriter{}
	rw.Read()
	rw.ReadWrite()
	rw.Write("hello")

}
