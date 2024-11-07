package main

import "go-learn/ch19/log"

func main() {

	log.Init(log.NewOptions())

	log.Debug("hello")
}
