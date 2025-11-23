package jepsen

import (
	"encoding/json"
	"fmt"
	"io"
)

type Node interface {
	fromInit(init initRequest, tx chan Event) Node
	processEvent(input Event, out io.Writer) error
}

func MainLoop(n Node, in io.Reader, out io.Writer) error {
	// Thread safe IO
	tx := make(chan Event)
	lw := &lockedWriter{Writer: out}
	lr := &lockedReader{Reader: in}
	enc := json.NewEncoder(lw)
	dec := json.NewDecoder(lr)

	// Read init message from stdin, else fail
	var init_msg message
	if err := dec.Decode(&init_msg); err != nil {
		return fmt.Errorf("failed to decode message: %w", err)
	}

	initBody, ok := init_msg.Body.(*initRequest)
	if !ok || initBody.Type != "init" {
		return fmt.Errorf("expected Init message, got invalid type")
	}

	// Send InitOk response
	reply := init_msg.intoReply(nil)
	if err := enc.Encode(reply); err != nil {
		return fmt.Errorf("failed to send init_ok: %w", err)
	}

	// Initialize Node and start processing events
	node := n.fromInit(*initBody, nil)
	go func() {
		for {
			var msg message
			if err := dec.Decode(&msg); err != nil {
				if err == io.EOF {
					tx <- &eof{}
					close(tx)
					return
				}
				// TODO:: Handle decode error
				continue
			}
			tx <- &msg
		}
	}()

	for event := range tx {
		if err := node.processEvent(event, lw); err != nil {
			return fmt.Errorf("failed to process payload: %w", err)
		}
	}

	return nil
}
