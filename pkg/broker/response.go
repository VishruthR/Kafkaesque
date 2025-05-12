package broker

import "fmt"

type response interface {
	getResponse() []byte
}

type pushResponse struct {
	status bool
}

func (r pushResponse) getResponse() []byte {
	if !r.status {
		return []byte("Push failed")
	}
	return []byte("Push succesful")
}

type pushNResponse struct {
	status    bool
	numPushed int
}

func (r pushNResponse) getResponse() []byte {
	if !r.status {
		return []byte("Push failed")

	}
	return []byte(fmt.Sprintf("Pushed %v elements succesfully", r.numPushed))
}

type pullResponse struct {
	status bool
	empty  bool
	body   string
}

func (r pullResponse) getResponse() []byte {
	if !r.status {
		return []byte("Push failed")
	}

	if r.empty {
		return []byte("Queue empty")
	}

	return []byte(fmt.Sprintf("%s\v", r.body))
}
