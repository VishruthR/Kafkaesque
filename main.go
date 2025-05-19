package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/VishruthR/Kafkaesque/pkg/broker"
)

func main() {
	recoverFlag := flag.Bool("recover", false, "When true, the broker will recover using the WAL")
	flag.Parse()
	args := flag.Args()
	allArgs := os.Args
	if len(args) != 1 {
		fmt.Printf("Invalid arguments\nusage:\t%s [--recover] port\n", allArgs[0])
		return
	}

	broker.Listen(args[0], *recoverFlag)
}
