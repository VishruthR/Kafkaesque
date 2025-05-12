package main

/*
Toy publisher client for testing purposes
*/

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

const TYPE = "tcp" // broker only supports tcp for now

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Printf("Invalid argumens\nUsage:\t<command> <host> <port>\n")
		return
	}

	host := args[1]
	port := args[2]

	reader := bufio.NewReader(os.Stdin)

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

	buf := make([]byte, 4096)
	for {
		numBytes, err := reader.Read(buf[:cap(buf)])
		buf = buf[:numBytes] // Select only newly read bytes
		if numBytes == 0 {
			if err == nil {
				continue
			} else if err == io.EOF {
				break
			}
			panic("Error reading STDIN")
		}

		_, err = conn.Write(buf)
		if err != nil {
			panic("Write to sock failed")
		}
		conn.CloseWrite()
	}

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
