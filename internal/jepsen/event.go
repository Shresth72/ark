package jepsen

import (
	"encoding/json"
	"io"
)

type MessageType string

const (
	TypeInit        MessageType = "init"
	TypeEcho        MessageType = "echo"
	TypeGenerate    MessageType = "generate"
	TypeAddCounter  MessageType = "add"
	TypeReadCounter MessageType = "read"
)

var messageBodyMap = map[MessageType]func() body{
	TypeInit:        func() body { return &initRequest{} },
	TypeEcho:        func() body { return &echoRequest{} },
	TypeGenerate:    func() body { return &uniqueIdsRequest{} },
	TypeAddCounter:  func() body { return &addToCounterRequest{} },
	TypeReadCounter: func() body { return &readCounterRequest{} },
}

// Event
type Event interface {
	isEvent()
	// intoReply(*int) Event
	// send(io.Writer) error
}

// Body
type body interface {
	isBody()
	intoReply(id int) body
}

// Message
type message struct {
	Src  string `json:"src"`
	Dst  string `json:"dest"`
	Body body   `json:"body"`
}

type injectMessage struct {
	Msg message
}

type eof struct{}

func (message) isEvent()       {}
func (injectMessage) isEvent() {}
func (eof) isEvent()           {}

func (m *message) intoReply(id *int) message {
	var newId int
	if id != nil {
		newId = *id
		*id++
	}

	return message{
		Src:  m.Dst,
		Dst:  m.Src,
		Body: m.Body.intoReply(newId),
	}
}

func (m *message) send(out io.Writer) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = out.Write(append(data, '\n'))
	if err != nil {
		return err
	}
	return nil
}

func (m *message) UnmarshalJSON(data []byte) error {
	var aux struct {
		Src  string          `json:"src"`
		Dst  string          `json:"dest"`
		Body json.RawMessage `json:"body"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	m.Src = aux.Src
	m.Dst = aux.Dst

	var bodyType struct {
		Type MessageType `json:"type"`
	}
	if err := json.Unmarshal(aux.Body, &bodyType); err != nil {
		return err
	}

	factory, ok := messageBodyMap[bodyType.Type]
	if !ok {
		return nil
	}

	m.Body = factory()
	return json.Unmarshal(aux.Body, m.Body)
}
