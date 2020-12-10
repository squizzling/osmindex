package state

import (
	"sync/atomic"
)

type Writer struct {
	currentState uint64
	currentBlock uint64
	copyAmount   uint64
	copyTotal    uint64
}

const (
	WriterReading = iota
	WriterWriting
)

func (w *Writer) CurrentState() uint64 {
	return atomic.LoadUint64(&w.currentState)
}

func (w *Writer) SetCurrentState(s uint64) {
	atomic.StoreUint64(&w.currentState, s)
}

func (w *Writer) SetStateWriting(block uint64) {
	atomic.StoreUint64(&w.currentBlock, block)
	atomic.StoreUint64(&w.currentState, WriterWriting)
}

func (w *Writer) CurrentBlock() uint64 {
	return atomic.LoadUint64(&w.currentBlock)
}
