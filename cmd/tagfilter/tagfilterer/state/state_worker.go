package state

import (
	"sync/atomic"
)

const (
	WorkerReading = iota
	WorkerDecoding
	WorkerWorking
	WorkerEncoding
	WorkerWriting
)

type Worker struct {
	currentState uint64
	currentBlock uint64
	kept         uint64
	dropped      uint64
}
func (w *Worker) CurrentState() uint64 {
	return atomic.LoadUint64(&w.currentState)
}

func (w *Worker) SetCurrentState(s uint64) {
	atomic.StoreUint64(&w.currentState, s)
}

func (w *Worker) CurrentBlock() uint64 {
	return atomic.LoadUint64(&w.currentBlock)
}
func (w *Worker) SetStateDecoding(block uint64) {
	atomic.StoreUint64(&w.currentBlock, block)
	atomic.StoreUint64(&w.currentState, WorkerDecoding)
}

func (w *Worker) AddKept(amount uint64) {
	atomic.AddUint64(&w.kept, amount)
}

func (w *Worker) AddDropped(amount uint64) {
	atomic.AddUint64(&w.dropped, amount)
}

func (w *Worker) Kept() uint64 {
	return atomic.LoadUint64(&w.kept)
}

func (w *Worker) Dropped() uint64 {
	return atomic.LoadUint64(&w.dropped)
}
