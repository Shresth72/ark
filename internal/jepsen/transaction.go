package jepsen

import (
	"fmt"
	"io"
	"sync"
)

type transactionRequest struct {
	Type string  `json:"type"`
	Txn  [][]any `json:"txn"`
}

type transactionResponse struct {
	Type      string  `json:"type"`
	MsgId     int     `json:"msg_id"`
	InReplyTo int     `json:"in_reply_to"`
	Txn       [][]any `json:"txn"`
}

func (transactionRequest) isBody()  {}
func (transactionResponse) isBody() {}

func (t transactionRequest) intoReply(id int) body {
	return &transactionResponse{
		Type:      "txn_ok",
		MsgId:     id,
		InReplyTo: id,
		Txn:       [][]any{},
	}
}

func (t transactionResponse) intoReply(id int) body {
	return &t
}

type TransactionNode struct {
	id   int
	node string

	mu    sync.Mutex
	store map[int]any
}

func (t *TransactionNode) fromInit(init initRequest, tx chan Event) Node {
	return &TransactionNode{
		id:    1,
		node:  init.NodeId,
		store: map[int]any{},
	}
}

func (t *TransactionNode) processEvent(input Event, out io.Writer) error {
	switch msg := input.(type) {
	case *message:
		switch body := msg.Body.(type) {
		case *transactionRequest:
			operation := t.handleTransaction(body.Txn)

			reply := msg.intoReply(&t.id)
			if r, ok := reply.Body.(*transactionResponse); ok {
				r.Txn = operation
			}
			return reply.send(out)

		default:
			return nil
		}
	default:
		return fmt.Errorf("Event type not allowed: %T", input)
	}
}

func (t *TransactionNode) handleTransaction(txns [][]any) [][]any {
	operation := make([][]any, 0)

	t.mu.Lock()
	defer t.mu.Unlock()

	for _, txn := range txns {
		if len(txn) != 3 {
			continue
		}
		op, key, value := txn[0], txn[1], txn[2]

		opStr, ok := op.(string)
		if !ok {
			continue
		}

		var keyInt int
		switch k := key.(type) {
		case int:
			keyInt = k
		case float64:
			keyInt = int(k)
		default:
			continue
		}

		switch opStr {
		case "r":
			val, ok := t.store[keyInt]
			if !ok {
				val = nil
			}
			operation = append(operation, []any{opStr, keyInt, val})
		case "w":
			t.store[keyInt] = value
		}
	}
	return operation
}
