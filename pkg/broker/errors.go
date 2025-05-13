package broker

import "fmt"

type InvalidElementInQueue struct {
	item interface{} // Catch all type, we want to just have the item so we can debug with it later
}

func (e *InvalidElementInQueue) Error() string {
	return fmt.Sprintf("Invalid element in broker queue %v\n", e.item)
}

type InvalidTopic struct {
	topic string
}

func (e *InvalidTopic) Error() string {
	return fmt.Sprintf("Invalid topic: %s\n", e.topic)
}

type TopicAlreadyExists struct {
	topic string
}

func (e *TopicAlreadyExists) Error() string {
	return fmt.Sprintf("Topic %s already exists\n", e.topic)
}
