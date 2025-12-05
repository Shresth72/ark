package jepsen

import (
	"fmt"
	"io"
	"strconv"
	"sync"
)

// Transaction DTO
type sendLogRequest struct {
	Type string `json:"type"`
	Key  string `json:"key"`
	Msg  int    `json:"msg"`
}

type sendLogResponse struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
}

type pollLogRequest struct {
	Type    string         `json:"type"`
	Offsets map[string]int `json:"offsets"`
}

type pollLogResponse struct {
	Type     string             `json:"type"`
	Messages map[string][][]int `json:"messages"`
}

type commitOffsetsRequest struct {
	Type    string         `json:"type"`
	Offsets map[string]int `json:"offsets"`
}

type commitOffsetsResponse struct {
	Type string `json:"type"`
}

type listCommitedOffsetsRequest struct {
	Type string   `json:"type"`
	Keys []string `json:"keys"`
}

type listCommitedOffsetsResponse struct {
	Type    string         `json:"type"`
	Offsets map[string]int `json:"offsets"`
}

func (sendLogRequest) isBody()              {}
func (sendLogResponse) isBody()             {}
func (pollLogRequest) isBody()              {}
func (pollLogResponse) isBody()             {}
func (commitOffsetsRequest) isBody()        {}
func (commitOffsetsResponse) isBody()       {}
func (listCommitedOffsetsRequest) isBody()  {}
func (listCommitedOffsetsResponse) isBody() {}

func (s sendLogRequest) intoReply(id int) body {
	return &sendLogResponse{
		Type:   "send_ok",
		Offset: 0,
	}
}

func (s sendLogResponse) intoReply(id int) body {
	return &s
}

func (p pollLogRequest) intoReply(id int) body {
	return &pollLogResponse{
		Type:     "poll_ok",
		Messages: map[string][][]int{},
	}
}

func (p pollLogResponse) intoReply(id int) body {
	return &p
}

func (c commitOffsetsRequest) intoReply(id int) body {
	return &commitOffsetsResponse{
		Type: "commit_offsets_ok",
	}
}

func (c commitOffsetsResponse) intoReply(id int) body {
	return &c
}

func (l listCommitedOffsetsRequest) intoReply(id int) body {
	return &listCommitedOffsetsResponse{
		Type:    "list_committed_offsets_ok",
		Offsets: map[string]int{},
	}
}

func (l listCommitedOffsetsResponse) intoReply(id int) body {
	return &l
}

// Transaction Node
type LogNode struct {
	id   int
	node string

	mu             sync.Mutex
	log_msgs       map[string][][]int
	committed_msgs map[string]int
}

func (l *LogNode) fromInit(init initRequest, tx chan Event) Node {
	return &LogNode{
		id:             1,
		node:           init.NodeId,
		log_msgs:       map[string][][]int{},
		committed_msgs: map[string]int{},
	}
}

func (l *LogNode) processEvent(input Event, out io.Writer) error {
	switch msg := input.(type) {
	case *message:
		switch body := msg.Body.(type) {
		case *sendLogRequest:
			offset, err := l.handleSendEvent(body.Key, body.Msg)
			if err != nil {
				return err
			}

			reply := msg.intoReply(&l.id)
			if r, ok := reply.Body.(*sendLogResponse); ok {
				r.Offset = offset
			}
			return reply.send(out)

		case *pollLogRequest:
			messages := l.handlePollEvent(body.Offsets)

			reply := msg.intoReply(&l.id)
			if r, ok := reply.Body.(*pollLogResponse); ok {
				r.Messages = messages
			}
			return reply.send(out)

		case *commitOffsetsRequest:
			l.mu.Lock()
			for key, offset := range body.Offsets {
				l.committed_msgs[key] = offset
			}
			l.mu.Unlock()

			reply := msg.intoReply(&l.id)
			return reply.send(out)

		case *listCommitedOffsetsRequest:
			l.mu.Lock()
			offsets := map[string]int{}
			for _, key := range body.Keys {
				offset, ok := l.committed_msgs[key]
				if ok {
					offsets[key] = offset
				}
			}
			l.mu.Unlock()

			reply := msg.intoReply(&l.id)
			if r, ok := reply.Body.(*listCommitedOffsetsResponse); ok {
				r.Offsets = offsets
			}
			return reply.send(out)

		default:
			return fmt.Errorf("unexpected message body type: %v", body)
		}
	default:
		return fmt.Errorf("Event type not allowed: %v", input)
	}
}

func (l *LogNode) handlePollEvent(offsets map[string]int) map[string][][]int {
	l.mu.Lock()
	defer l.mu.Unlock()

	result := map[string][][]int{}

	for key, vecs := range l.log_msgs {
		startOffset := offsets[key]
		startOffset = min(vecs[0][0], startOffset)

		filtered := [][]int{}
		for _, vec := range vecs {
			if vec[0] >= startOffset {
				filtered = append(filtered, vec)
			}
		}
		if len(filtered) > 0 {
			result[key] = filtered
		}
	}

	return result
}

func (l *LogNode) handleSendEvent(key string, msg int) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.log_msgs[key]
	if !ok {
		entry = [][]int{}
	}

	var newOffset int
	if len(entry) == 0 {
		idx := -1
		for i, r := range key {
			if r >= '0' && r <= '9' {
				idx = i
				break
			}
		}
		if idx == -1 {
			return 0, fmt.Errorf("Invalid message key: %s", key)
		}

		num, err := strconv.Atoi(key[idx:])
		if err != nil {
			return 0, fmt.Errorf("Invalid message key: %s", key)
		}
		newOffset = num * 1000
	} else {
		last := entry[len(entry)-1]
		maxOffset := last[0]
		newOffset = maxOffset + 1
	}

	entry = append(entry, []int{newOffset, msg})
	l.log_msgs[key] = entry

	return newOffset, nil
}
