package state

import (
	"sync/atomic"
)

type Writer struct {
	currentState uint64
	currentBlock uint64
	fillerBlocks uint64
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

func (w *Writer) AddFiller() {
	atomic.AddUint64(&w.fillerBlocks, 1)
}

func (w *Writer) Filler() uint64 {
	return atomic.LoadUint64(&w.fillerBlocks)
}
