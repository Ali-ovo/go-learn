package main

import (
	"flag"
	"fmt"
)

func main() {

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Useage: admin-cli [command]")
		return

	}

	switch args[0] {
	case "help":
		fmt.Println("help")
	case "export":
		if len(args) == 2 {
			fmt.Println("export to file")
		} else if len(args) == 1 {
			fmt.Println("export to  default file")
		}
	default:
		fmt.Println("Unknown command")
	}
}
