// WAL writes must be synchronous
// Checkpointing can be async
package broker

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func AppendLog(walWriter *bufio.Writer, request *request) error {
	log := fmt.Sprintf("%s;%s;%v %s\n", request.header.request.String(), request.header.topic, request.header.contentLength, request.body)
	if _, err := walWriter.WriteString(log); err != nil {
		return err
	}

	if err := walWriter.Flush(); err != nil {
		return err
	}

	return nil
}

func logToRequest(log string) (*request, error) {
	logParts := strings.Split(strings.TrimSpace(log), " ")

	headerParts := strings.Split(logParts[0], ";")
	if len(headerParts) != 3 {
		return nil, &InvalidMessage{"Invalid header format"}
	}

	requestType, err := stringToRequestType(headerParts[0])
	if err != nil {
		return nil, err
	}

	contentLen, err := strconv.ParseUint(headerParts[2], 10, 64)
	if err != nil {
		return nil, &InvalidMessage{"Invalid content length"}
	}

	body := ""
	if len(logParts) == 2 {
		body = logParts[1]
	}

	return &request{
		requestHeader{
			requestType,
			headerParts[1],
			contentLen,
		},
		body,
	}, nil
}

func RecoverFromLog(topics *topics) error {
	wal, err := os.Open("wal.log")
	if err != nil {
		return err
	}
	defer wal.Close()

	scanner := bufio.NewScanner(wal)
	for scanner.Scan() {
		request, err := logToRequest(scanner.Text())
		if err != nil {
			return err
		}

		_, err = request.processRequest(topics)
		if err != nil {
			fmt.Println("Error processing a request from WAL")
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
