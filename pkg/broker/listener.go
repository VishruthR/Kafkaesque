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

func Listen(port string) {
	broker := brokerQueue{
		q: queue.New(),
	}
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
		go handleConnection(conn, &broker)
	}
}

func processPacket(readBuf []byte, writeBuf *[]byte) (*request, error) {
	// fmt.Println(string(readBuf))
	messageParts := strings.Split(string(readBuf), HEADER_END)
	if len(messageParts) != 2 {
		panic("Invalid message structure; Check if only header and body are presents") // TODO: processPacket should return error here instead of panicking
	}

	header := messageParts[0]
	header, exists := strings.CutPrefix(header, HEADER_START)
	if !exists {
		// While it isn't strictly necessary to panic, it is probably good to enforce strict schema
		panic("Invalid message structure; Ensure header has both opening and closing tag")
	}
	body := messageParts[1]
	fmt.Printf("%s\n", body)
	body, existsPre := strings.CutPrefix(body, BODY_START)
	body, existsSuf := strings.CutSuffix(body, BODY_END)
	if !(existsPre && existsSuf) {
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

	request := request{
		requestHeader{
			stringToRequestType(headerParts[0]),
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

func handleConnection(conn net.Conn, broker *brokerQueue) {
	defer conn.Close()
	defer recoverFromParsingError(conn)
	writeBuf := make([]byte, 4096)
	readBuf := make([]byte, 0, 4096)
	for {
		temp := make([]byte, 4096)
		bytesRead, err := conn.Read(temp)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Connection %s read error: %v\n", conn.RemoteAddr().String(), err) // TODO: Better error handling
			}
			break
		}

		finalLen := len(readBuf) + bytesRead
		readBuf = append(readBuf, temp...)
		readBuf = readBuf[:finalLen]
		if isFullPacket(readBuf) {
			request, err := processPacket(readBuf, &writeBuf)
			if err != nil {
				panic("Error processing packet") // TOOD: Fix error handling of processPacket
			}
			fmt.Printf("%v\n", request)

			response, err := request.processRequest(broker)
			if err != nil {
				panic("Error processing request!")
			}

			writeBuf = append(writeBuf, response.getResponse()...)
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
