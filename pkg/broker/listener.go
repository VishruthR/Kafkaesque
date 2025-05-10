package broker

import (
	"fmt"
	"io"
	"log"
	"net"
)

func Listen(port string) {
	listener, err := net.Listen("tcp4", ":" + port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Printf("Listening on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err) // TODO: Better error handling; ideally a failed connection attempt should not crash the broker
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Connected to %s\n", conn.RemoteAddr().String())
	writeBuf := make([]byte, 4096)
	readBuf := make([]byte, 4096)
	for {
		_, err := conn.Read(readBuf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Connection %s read error: %v\n", conn.RemoteAddr().String(), err) // TODO: Better error handling
			}
			break
		}
		writeBuf = append(writeBuf, readBuf...)
	}
	conn.Write(writeBuf)
}