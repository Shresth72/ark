package jepsen

import (
	"fmt"
	"io"
)

// Echo DTO
type echoRequest struct {
	Type string `json:"type"`
	Id   int    `json:"msg_id"`
	Echo string `json:"echo"`
}

type echoOkResponse struct {
	Type      string `json:"type"`
	Id        int    `json:"msg_id"`
	InReplyTo int    `json:"in_reply_to"`
	Echo      string `json:"echo"`
}

func (echoRequest) isBody()    {}
func (echoOkResponse) isBody() {}

func (e echoRequest) intoReply(id int) body {
	return &echoOkResponse{
		Type:      "echo_ok",
		Id:        id,
		InReplyTo: e.Id,
		Echo:      e.Echo,
	}
}

func (e echoOkResponse) intoReply(id int) body {
	return &e
}

// Echo Node
type EchoNode struct {
	id int
}

func (e *EchoNode) fromInit(init initRequest, tx chan Event) Node {
	return &EchoNode{id: 1}
}

func (e *EchoNode) processEvent(input Event, out io.Writer) error {
	switch msg := input.(type) {
	case *message:
		switch msg.Body.(type) {
		case *echoRequest:
			reply := msg.intoReply(&e.id)
			return reply.send(out)
		default:
			return nil
		}
	default:
		return fmt.Errorf("Event type not allowed: %T", input)
	}
}
