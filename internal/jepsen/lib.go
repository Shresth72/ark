package jepsen

import "io"

type Node interface {
	fromInit(init initRequest, tx chan Event) Node
	processEvent(input Event, out io.Writer) error
}
