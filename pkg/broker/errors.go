package broker

import "fmt"

type InvalidMessage struct {
	message string
}

func (e *InvalidMessage) Error() string {
	return fmt.Sprintf("Invalid message: %s", e.message)
}

type InvalidRequestType struct {
	requestString string
}

func (e *InvalidRequestType) Error() string {
	return fmt.Sprintf("Invalid request type: %s", e.requestString)
}

type InvalidElementInQueue struct {
	item interface{} // Catch all type, we want to just have the item so we can debug with it later
}

func (e *InvalidElementInQueue) Error() string {
	return fmt.Sprintf("Invalid element in broker queue %v", e.item)
}

type InvalidTopic struct {
	topic string
}

func (e *InvalidTopic) Error() string {
	return fmt.Sprintf("Invalid topic: %s", e.topic)
}

type TopicAlreadyExists struct {
	topic string
}

func (e *TopicAlreadyExists) Error() string {
	return fmt.Sprintf("Topic %s already exists", e.topic)
}
