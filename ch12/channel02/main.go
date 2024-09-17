package main

func main() {
	// 默认 channel 是双向的
	// var ch1 chan int       // 双向
	// var ch2 chan<- float64 // 单向写入
	// var ch3 <-chan int     // 单向读取

	c := make(chan int, 3)
	var send chan<- int = c // send only
	var read <-chan int = c // read only

	send <- 1

	<-read
}
