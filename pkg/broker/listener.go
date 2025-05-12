package broker

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	HEADER_START          = "<header>\n"
	HEADER_END            = "</header>\n"
	BODY_START            = "<body>\n"
	BODY_END              = "</body>"
	PULL_REQUEST_PROTOCOL = "PULL"
	PUSH_REQUEST_PROTOCOL = "PUSH"
)

// De facto enum
type RequestType int

const (
	PULL RequestType = iota
	PULL_N
	PULL_HEAD
	PUSH
	PUSH_N
)

func (r RequestType) String() string {
	switch r {
	case PULL:
		return "PULL"
	case PULL_N:
		return "PULL_N"
	case PULL_HEAD:
		return "PULL_HEAD"
	case PUSH:
		return "PUSH"
	case PUSH_N:
		return "PUSH_N"
	default:
		panic("Invalid Request")
	}
}

func stringToRequestType(s string) RequestType {
	switch s {
	case "PULL":
		return PULL
	case "PULL_N":
		return PULL_N
	case "PULL_HEAD":
		return PULL_HEAD
	case "PUSH":
		return PUSH
	case "PUSH_N":
		return PUSH_N
	default:
		panic("Invalid request string")
	}
}

type requestHeader struct {
	request       RequestType
	topic         string
	contentLength uint64
}

func (rH requestHeader) String() string {
	return fmt.Sprintf("{\n\trequest: %v,\n\ttopic: %v,\n\tcontentLength: %v\n}", rH.request, rH.topic, rH.contentLength)
}

func Listen(port string) {
	listener, err := net.Listen("tcp4", ":"+port)
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

func processPacket(readBuf []byte, writeBuf *[]byte) {
	// fmt.Println(string(readBuf))
	messageParts := strings.Split(string(readBuf), HEADER_END)
	if len(messageParts) != 2 {
		panic("Invalid message structure; Check if only header and body are presents")
	}

	header := messageParts[0]
	header, exists := strings.CutPrefix(header, HEADER_START)
	if !exists {
		// While it isn't strictly necessary to panic, it is probably good to enforce strict schema
		panic("Invalid message structure; Ensure header has both opening and closing tag")
	}
	body := messageParts[1]
	body, exists = strings.CutPrefix(body, BODY_START)
	if !exists {
		panic("Invalid message structure; Ensure body has both opening and closing tag")
	}

	headerParts := strings.SplitN(header, "\n", 3)
	if len(headerParts) != 3 {
		panic("Invalid header format")
	}
	contentLen, err := strconv.ParseUint(headerParts[2][:len(headerParts[2])-1], 10, 64)
	if err != nil {
		panic("Invalid content length")
	}

	requestHeader := requestHeader{
		stringToRequestType(headerParts[0]),
		headerParts[1],
		contentLen,
	}
	fmt.Println(requestHeader)

	(*writeBuf) = append((*writeBuf), body...)
}

func isFullPacket(packet []byte) bool {
	return bytes.Contains(packet, []byte(BODY_END)) // TODO: Make this more efficient/comprehensive
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	defer recoverFromParsingError(conn)
	writeBuf := make([]byte, 4096)
	readBuf := make([]byte, 0, 4096)
	for {
		temp := make([]byte, 4096)
		_, err := conn.Read(temp)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Connection %s read error: %v\n", conn.RemoteAddr().String(), err) // TODO: Better error handling
			}
			break
		}

		readBuf = append(readBuf, temp...)
		if isFullPacket(readBuf) {
			processPacket(readBuf, &writeBuf)
			conn.Write(writeBuf)
			break
		}

	}
	conn.Write(writeBuf)
}

func recoverFromParsingError(conn net.Conn) {
	if r := recover(); r != nil {
		fmt.Printf("Connection %s error: %v\n", conn.RemoteAddr().String(), r)
		conn.Write([]byte(fmt.Sprintf("%v", r)))
	}
}
