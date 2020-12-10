package state

import (
	"sync/atomic"
)

const (
	WorkerReading = iota
	WorkerDecoding
	WorkerWorking
	WorkerWriting
)

type Worker struct {
	currentState uint64
	currentBlock uint64
	currentWay   uint64
	wayCount     uint64
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

func (w *Worker) SetStateWorking(wayCount int) {
	atomic.StoreUint64(&w.wayCount, uint64(wayCount))
	atomic.StoreUint64(&w.currentWay, 0)
	atomic.StoreUint64(&w.currentState, WorkerWorking)
}

func (w *Worker) NextWay() {
	atomic.AddUint64(&w.currentWay, 1)
}

func (w *Worker) CurrentWay() uint64 {
	return atomic.LoadUint64(&w.currentWay)
}

func (w *Worker) WayCount() uint64 {
	return atomic.LoadUint64(&w.wayCount)
}
