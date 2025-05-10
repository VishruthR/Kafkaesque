package main

import (
	"os"
	"fmt"
	"github.com/VishruthR/Kafkaesque/pkg/broker"
)

func main () {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Invalid arguments\nUsage:\t<command> <port>\n")
		return
	}

	broker.Listen(args[1])
}
