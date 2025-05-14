package broker

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/phf/go-queue/queue"
)

type brokerQueue struct {
	q  *queue.Queue
	mu sync.Mutex
}

type topics map[string]*brokerQueue

func Listen(port string) {
	topics := make(topics)
	listener, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Printf("Listening on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue // Skip this connection
		}
		go handleConnection(conn, &topics)
	}
}

func processPacket(readBuf []byte) (*request, error) {
	messageParts := strings.Split(string(readBuf), HEADER_END)
	if len(messageParts) != 2 {
		return nil, &InvalidMessage{"Message must only contain header and body sections"}
	}

	header := messageParts[0]
	header, exists := strings.CutPrefix(header, HEADER_START)
	if !exists {
		return nil, &InvalidMessage{"Ensure header has both opening and closing tag"}
	}
	body := messageParts[1]
	fmt.Printf("%s\n", body)
	body, existsPre := strings.CutPrefix(body, BODY_START)
	body, existsSuf := strings.CutSuffix(body, BODY_END)
	if !(existsPre && existsSuf) {
		return nil, &InvalidMessage{"Ensure body has both opening and closing tag"}
	}

	headerParts := strings.SplitN(header, "\n", 3)
	if len(headerParts) != 3 {
		return nil, &InvalidMessage{"Invalid header format"}
	}
	headerParts[2] = strings.TrimSpace(headerParts[2])
	contentLen, err := strconv.ParseUint(headerParts[2], 10, 64)
	if err != nil {
		return nil, &InvalidMessage{"Invalid content length"}
	}

	requestType, err := stringToRequestType(headerParts[0])
	if err != nil {
		return nil, err
	}

	request := request{
		requestHeader{
			requestType,
			headerParts[1],
			contentLen,
		},
		body, // copying whole string here could be a little slow
	}

	return &request, nil
}

func isFullPacket(packet []byte) bool {
	return bytes.Contains(packet, []byte(BODY_END)) // TODO: Make this more efficient/comprehensive
}

func handleConnection(conn net.Conn, topics *topics) {
	defer conn.Close()
	writeBuf := make([]byte, 4096)
	readBuf := make([]byte, 0, 4096)
	for {
		temp := make([]byte, 4096)
		bytesRead, err := conn.Read(temp)
		if err != nil {
			if err != io.EOF {
				handleConnectionError(conn, err)
				return
			}
			break
		}

		finalLen := len(readBuf) + bytesRead
		readBuf = append(readBuf, temp...)
		readBuf = readBuf[:finalLen]
		if isFullPacket(readBuf) {
			request, err := processPacket(readBuf)
			if err != nil {
				handleConnectionError(conn, err)
				return
			}
			fmt.Printf("%v\n", request)

			response, err := request.processRequest(topics)
			if err != nil {
				handleConnectionError(conn, err)
				return
			}

			writeBuf = append(writeBuf, response.getResponse()...)
			break
		}

	}

	conn.Write(writeBuf)
}

func handleConnectionError(conn net.Conn, e error) {
	fmt.Printf("Connection %s error: %v\n", conn.RemoteAddr().String(), e)
	conn.Write([]byte(fmt.Sprintf("%v", e)))
}
