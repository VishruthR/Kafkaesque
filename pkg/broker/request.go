package broker

import (
	"fmt"
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

type request struct {
	header requestHeader
	body   string
}

func (r request) String() string {
	return fmt.Sprintf("{\n\theader: %v,\n\tbody: %v\n}", r.header, r.body)
}

func (request request) processRequest(broker *brokerQueue) (response, error) {
	switch request.header.request {
	case PULL:
		return request.processPull(broker)
	case PULL_N:
		fmt.Println("Handling PULL_N")
	case PULL_HEAD:
		fmt.Println("Handling PULL_HEAD")
	case PUSH:
		return request.processPush(broker)
	case PUSH_N:
		return request.processPushN(broker)
	default:
		panic("Invalid request in requestHeader") // Panic because this should never happen
	}

	return pushResponse{status: false}, nil
}

func (request request) processPush(broker *brokerQueue) (pushResponse, error) {
	broker.mu.Lock()
	defer broker.mu.Unlock()
	broker.q.PushBack(request.body)

	return pushResponse{status: true}, nil
}

func (request request) processPushN(broker *brokerQueue) (pushNResponse, error) {
	items := strings.Split(request.body, ";")
	broker.mu.Lock()
	defer broker.mu.Unlock()
	for _, item := range items {
		broker.q.PushBack(item)
	}

	return pushNResponse{status: true, numPushed: len(items)}, nil
}

func (request request) processPull(broker *brokerQueue) (pullResponse, error) {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	item := broker.q.PopFront()
	if item == nil {
		return pullResponse{
			status: true,
			empty:  true,
		}, nil
	}

	item_str, ok := item.(string)
	if !ok {
		fmt.Println("Cannot convert item from queue to string")
		return pullResponse{
			status: false,
		}, nil // Should throw an actual error here
	}

	return pullResponse{
		status: true,
		empty:  false,
		body:   item_str,
	}, nil
}
