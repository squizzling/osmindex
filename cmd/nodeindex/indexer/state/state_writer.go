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
	WriterWritingIndex
	WriterWritingLocations
	WriterCopyingLocations
	WriteFinalizing
)

func (w *Writer) CurrentState() uint64 {
	return atomic.LoadUint64(&w.currentState)
}

func (w *Writer) SetCurrentState(s uint64) {
	atomic.StoreUint64(&w.currentState, s)
}

func (w *Writer) SetStateWritingIndex(block uint64) {
	atomic.StoreUint64(&w.currentBlock, block)
	atomic.StoreUint64(&w.currentState, WriterWritingIndex)
}

func (w *Writer) SetStateCopyLocations(sz uint64) {
	atomic.StoreUint64(&w.copyTotal, sz)
	atomic.StoreUint64(&w.currentState, WriterCopyingLocations)
}

func (w *Writer) CopyAmount() uint64 {
	return atomic.LoadUint64(&w.copyAmount)
}

func (w *Writer) AdvanceCopy(amount uint64) {
	atomic.AddUint64(&w.copyAmount, amount)
}

func (w *Writer) CopyTotal() uint64 {
	return atomic.LoadUint64(&w.copyTotal)
}

func (w *Writer) CurrentBlock() uint64 {
	return atomic.LoadUint64(&w.currentBlock)
}
