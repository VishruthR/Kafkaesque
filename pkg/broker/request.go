package broker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/phf/go-queue/queue"
)

const (
	HEADER_START = "<header>\n"
	HEADER_END   = "</header>\n"
	BODY_START   = "<body>\n"
	BODY_END     = "</body>"
)

// De facto enum
type RequestType int

const (
	PULL RequestType = iota
	PULL_N
	PULL_HEAD
	PUSH
	PUSH_N
	CREATE_TOPIC
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
	case CREATE_TOPIC:
		return "CREATE_TOPIC"
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
	case "CREATE_TOPIC":
		return CREATE_TOPIC
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

func (request *request) processRequest(topics *topics) (response, error) {
	switch request.header.request { // Handle non-topic specific requests
	case CREATE_TOPIC:
		return request.processCreateTopic(topics)
	}

	broker, exists := (*topics)[request.header.topic]
	if !exists {
		return invalidTopicResponse{
			topic: request.header.topic,
		}, &InvalidTopic{topic: request.header.topic}
	}

	switch request.header.request {
	case PULL:
		return request.processPull(broker)
	case PULL_N:
		return request.processPullN(broker)
	case PULL_HEAD:
		return request.processPullHead(broker)
	case PUSH:
		return request.processPush(broker)
	case PUSH_N:
		return request.processPushN(broker)
	default:
		panic("Invalid request in requestHeader") // Panic because this should never happen
	}
}

func (request *request) processPush(broker *brokerQueue) (pushResponse, error) {
	broker.mu.Lock()
	defer broker.mu.Unlock()
	broker.q.PushBack(request.body)

	return pushResponse{status: true}, nil
}

func (request *request) processPushN(broker *brokerQueue) (pushNResponse, error) {
	items := strings.Split(request.body, ";")
	broker.mu.Lock()
	defer broker.mu.Unlock()
	for _, item := range items {
		broker.q.PushBack(item)
	}

	return pushNResponse{status: true, numPushed: len(items)}, nil
}

func (request *request) processPull(broker *brokerQueue) (pullResponse, error) {
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
		}, &InvalidElementInQueue{item: item}
	}

	return pullResponse{
		status: true,
		empty:  false,
		body:   item_str,
	}, nil
}

func (request *request) processPullN(broker *brokerQueue) (pullNResponse, error) {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	request.body = strings.TrimSpace(request.body)
	numToPull, err := strconv.ParseUint(request.body, 10, 64)
	if err != nil {
		return pullNResponse{
			status: false,
		}, err
	}

	response := pullNResponse{
		body:      make([]string, 0, 64), // Initialize with some base capacity to avoid unecessary resizing
		numPulled: 0,
	}

	for response.numPulled < numToPull {
		item := broker.q.PopFront()
		if item == nil { // Reached end of queue
			break
		}

		item_str, ok := item.(string)
		if !ok {
			fmt.Println("Cannot convert item from queue to string")
			response.status = true // Not sure if I should return error or return the elements that we fetched
			return response, &InvalidElementInQueue{item: item}
		}

		response.body = append(response.body, item_str)
		response.numPulled += 1
	}

	response.status = true
	return response, nil
}

// This function lwk identical to pull N right now, but will keep separate for now in case of changes to API
func (request *request) processPullHead(broker *brokerQueue) (pullNResponse, error) {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	response := pullNResponse{
		body:      make([]string, 0, 64), // Initialize with some base capacity to avoid unecessary resizing
		numPulled: 0,
	}

	for {
		item := broker.q.PopFront()
		if item == nil { // Reached end of queue
			break
		}

		item_str, ok := item.(string)
		if !ok {
			fmt.Println("Cannot convert item from queue to string")
			response.status = true // Not sure if I should return error or return the elements that we fetched
			return response, &InvalidElementInQueue{item: item}
		}

		response.body = append(response.body, item_str)
		response.numPulled += 1
	}

	response.status = true
	return response, nil
}

func (request request) processCreateTopic(topics *topics) (createTopicResponse, error) {
	if _, exists := (*topics)[request.header.topic]; exists {
		return createTopicResponse{
			status: false, topic: request.header.topic,
		}, &TopicAlreadyExists{topic: request.header.topic}
	}

	newBroker := brokerQueue{
		q: queue.New(),
	}
	(*topics)[request.header.topic] = &newBroker

	return createTopicResponse{
		status: true, topic: request.header.topic,
	}, nil
}
