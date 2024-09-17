package main

import "time"

func main() {

	msg := make(chan string, 2)

	go func(msg chan string) {
		for data := range msg {
			println(data)
		}

		println("all done")
	}(msg)

	msg <- "1"
	msg <- "2"

	close(msg)

	time.Sleep(time.Second * 2)
}
