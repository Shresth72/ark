package jepsen

import (
	"fmt"
	"io"
	"sync"
)

// Grow Counter DTO
type addToCounterRequest struct {
	Type  string `json:"type"`
	Delta int    `json:"delta"`
}

type addToCounterResponse struct {
	Type string `json:"type"`
}

type readCounterRequest struct {
	Type string `json:"type"`
}

type readCounterResponse struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

func (addToCounterRequest) isBody()  {}
func (addToCounterResponse) isBody() {}
func (readCounterRequest) isBody()   {}
func (readCounterResponse) isBody()  {}

func (a addToCounterRequest) intoReply(id int) body {
	return &addToCounterResponse{
		Type: "add_ok",
	}
}

func (a addToCounterResponse) intoReply(id int) body {
	return &a
}

func (r readCounterRequest) intoReply(id int) body {
	return &readCounterResponse{
		Type:  "read_ok",
		Value: 0,
	}
}

func (r readCounterResponse) intoReply(id int) body {
	return &r
}

// Grow Counter Node
type GrowCounterNode struct {
	id   int
	node string

	mu      sync.Mutex
	counter int
}

func (g *GrowCounterNode) fromInit(init initRequest, tx chan Event) Node {
	return &GrowCounterNode{
		id:      1,
		node:    init.NodeId,
		counter: 0,
	}
}

func (g *GrowCounterNode) processEvent(input Event, out io.Writer) error {
	switch msg := input.(type) {
	case *message:
		switch body := msg.Body.(type) {
		case *addToCounterRequest:
			delta := body.Delta

			g.mu.Lock()
			g.counter += delta
			g.mu.Unlock()

			reply := msg.intoReply(&g.id)
			return reply.send(out)

		case *readCounterRequest:
			g.mu.Lock()
			value := g.counter
			g.mu.Unlock()

			reply := msg.intoReply(&g.id)
			if r, ok := reply.Body.(*readCounterResponse); ok {
				r.Value = value
			}
			return reply.send(out)

		default:
			return fmt.Errorf("unexpected message body type: %v", msg.Body)
		}
	default:
		return fmt.Errorf("Event type not allowed: %v", input)
	}
}
