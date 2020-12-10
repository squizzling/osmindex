package state

import (
	"sync/atomic"
)

const (
	WorkerReceivingBlock = iota
	WorkerDecodingBlob
	WorkerDecompressingBlob
	WorkerEncodingBlob
	WorkerSendingBlock
)

type Worker struct {
	currentState uint64
	currentBlock uint64
}

func (sw *Worker) CurrentState() uint64 {
	return atomic.LoadUint64(&sw.currentState)
}

func (sw *Worker) SetCurrentState(s uint64) {
	atomic.StoreUint64(&sw.currentState, s)
}

func (sw *Worker) SetDecodingBlob(block uint64) {
	atomic.StoreUint64(&sw.currentBlock, block)
	atomic.StoreUint64(&sw.currentState, WorkerDecodingBlob)
}

func (sw *Worker) CurrentBlock() uint64 {
	return atomic.LoadUint64(&sw.currentBlock)
}
