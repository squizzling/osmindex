package state

import (
	"sync/atomic"
)

const (
	WorkerReading = iota
	WorkerDecoding
	WorkerProcessing
	WorkerWriting
)

type Worker struct {
	currentState uint64
	currentBlock uint64
}

func (w *Worker) CurrentState() uint64 {
	return atomic.LoadUint64(&w.currentState)
}

func (w *Worker) SetCurrentState(s uint64) {
	atomic.StoreUint64(&w.currentState, s)
}

func (w *Worker) SetStateDecodeBlock(block uint64) {
	atomic.StoreUint64(&w.currentBlock, block)
	atomic.StoreUint64(&w.currentState, WorkerDecoding)
}

func (w *Worker) CurrentBlock() uint64 {
	return atomic.LoadUint64(&w.currentBlock)
}
