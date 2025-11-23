package jepsen

import (
	"io"
	"sync"
)

type lockedWriter struct {
	io.Writer
	sync.Mutex
}

func (lw *lockedWriter) Write(p []byte) (n int, err error) {
	lw.Lock()
	defer lw.Unlock()
	return lw.Writer.Write(p)
}

type lockedReader struct {
	io.Reader
	sync.Mutex
}

func (lw *lockedReader) Read(p []byte) (n int, err error) {
	lw.Lock()
	defer lw.Unlock()
	return lw.Reader.Read(p)
}
