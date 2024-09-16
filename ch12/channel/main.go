package main

func main() {
	// var msg chan string
	// msg = make(chan string, 0)

	msg := make(chan string, 0)

	go func(msg chan string) {
		s := <-msg
		println(s)
	}(msg)

	msg <- "hello"

}
