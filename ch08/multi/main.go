package main

import "fmt"

type MyWriter interface {
	Write(string) error
}

type MyCloser interface {
	Close() error
}

type writerCloser struct {
	MyWriter
}

type fileWriter struct {
	filePath string
}

type databaseWriter struct {
	host string
}

func (wc *fileWriter) Write(s string) error {
	fmt.Println("write string to file")
	return nil
}

func (wc *databaseWriter) Write(s string) error {
	fmt.Println("write string to database")
	return nil
}

func (wc *writerCloser) Close() error {
	fmt.Println("close string")
	return nil
}

func main() {
	// var mw MyWriter = &writerCloser{}
	// mw.Write("hello")

	// var mc MyCloser = &writerCloser{}
	// mc.Close()

	var mw MyWriter = &writerCloser{
		&fileWriter{
			filePath: "a.txt",
		},
	}

	var mc MyWriter = &writerCloser{
		&databaseWriter{
			host: "localhost",
		},
	}

	mw.Write("hello")
	mc.Write("hello")

}
