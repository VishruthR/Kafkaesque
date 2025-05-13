package broker

import (
	"fmt"
	"strings"
)

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
		return []byte("Pull failed")
	}

	if r.empty {
		return []byte("Queue empty")
	}

	return []byte(fmt.Sprintf("%s\v", r.body))
}

type pullNResponse struct {
	status    bool
	numPulled uint64
	body      []string
}

func (r pullNResponse) getResponse() []byte {
	if !r.status {
		return []byte("Pull failed")
	}

	if r.numPulled == 0 {
		return []byte("Queue empty")
	}

	return []byte(fmt.Sprintf("%v;%v", r.numPulled, strings.Join(r.body, ";")))
}
