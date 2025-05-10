package main

/*
Toy publisher client for testing purposes
*/

import (
	"fmt"
	"io"
	"net"
	"os"
)

const TYPE = "tcp" // broker only supports tcp for now

func main() {
	args := os.Args
	if len(args) != 4 {
		fmt.Printf("Invalid argumens\nUsage:\t<command> <host> <port> <payload>\n")
		return
	}

	host := args[1]
	port := args[2]
	payload := args[3]

	tcpServer, err := net.ResolveTCPAddr(TYPE, host+":"+port)
	if err != nil {
		fmt.Println("ResolveTCPAddr failed:", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		fmt.Println("Unable to connect:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(payload))
	if err != nil {
		fmt.Println("Write failed:", err)
		return
	}
	conn.CloseWrite()

	received := make([]byte, 4096)
	for {
		readBuf := make([]byte, 4096)
		_, err = conn.Read(readBuf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Read failed:", err)
			return
		}
		received = append(received, readBuf...)
	}

	fmt.Printf("Receieved: %s\n", string(received))
}