package state

import (
	"sync/atomic"
)

const (
	WriterReceivingBlock = iota
	WriterWritingData
)

type Writer struct {
	currentState uint64
	currentBlock uint64
}

func (d *Writer) CurrentState() uint64 {
	return atomic.LoadUint64(&d.currentState)
}

func (d *Writer) SetCurrentState(s uint64) {
	atomic.StoreUint64(&d.currentState, s)
}

func (d *Writer) SetWritingBlock(block uint64) {
	atomic.StoreUint64(&d.currentBlock, block)
	atomic.StoreUint64(&d.currentState, WriterWritingData)
}

func (d *Writer) CurrentBlock() uint64 {
	return atomic.LoadUint64(&d.currentBlock)
}
