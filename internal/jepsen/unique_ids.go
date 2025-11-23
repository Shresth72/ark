package jepsen

import (
	"fmt"
	"io"
)

// UniqueIds DTO
type uniqueIdsRequest struct {
	Type string `json:"type"`
}

type uniqueIdsResponse struct {
	Type string `json:"type"`
	Id   int    `json:"id"`
}

func (uniqueIdsRequest) isBody()  {}
func (uniqueIdsResponse) isBody() {}

func (u uniqueIdsRequest) intoReply(id int) body {
	return &uniqueIdsResponse{
		Type: "generate",
		Id:   id,
	}
}

func (u uniqueIdsResponse) intoReply(id int) body {
	return &u
}

// UniqueIds Node
type UniqueIdsNode struct {
	id   int
	node string
}

func (u *UniqueIdsNode) fromInit(init initRequest, tx chan Event) Node {
	return &UniqueIdsNode{
		id:   1,
		node: init.NodeId,
	}
}

func (u *UniqueIdsNode) processEvent(input Event, out io.Writer) error {
	switch msg := input.(type) {
	case *message:
		switch msg.Body.(type) {
		case echoRequest:
			reply := msg.intoReply(&u.id)
			return reply.send(out)
		default:
			return nil
		}
	default:
		return fmt.Errorf("Event type not allowed: %T", input)
	}
}
