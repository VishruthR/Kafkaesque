package broker

import "fmt"

type InvalidElementInQueue struct {
	item interface{} // Catch all type, we want to just have the item so we can debug with it later
}

func (e *InvalidElementInQueue) Error() string {
	return fmt.Sprintf("Invalid element in broker queue %v\n", e.item)
}
