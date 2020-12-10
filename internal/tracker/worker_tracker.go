package tracker

import (
	"sync"
	"sync/atomic"
)

// WorkerTracker maintains a list of worker states
type WorkerTracker struct {
	version uint64 // atomic, mu does not protect this
	mu      sync.RWMutex
	workers []interface{}
}

func (wt *WorkerTracker) Version() uint64 {
	return atomic.LoadUint64(&wt.version)
}

func (wt *WorkerTracker) TrackWorker(sw interface{}) {
	wt.mu.Lock()
	wt.workers = append(wt.workers, sw)
	wt.mu.Unlock()
	atomic.AddUint64(&wt.version, 1)
}

func (wt *WorkerTracker) UntrackWorker(sw interface{}) {
	wt.mu.Lock()
	for idx, swOther := range wt.workers {
		if swOther == sw {
			copy(wt.workers[idx:], wt.workers[idx+1:])
			wt.workers[len(wt.workers)-1] = nil
			wt.workers = wt.workers[:len(wt.workers)-1]
			break
		}
	}
	wt.mu.Unlock()
	atomic.AddUint64(&wt.version, 1)
}

func (wt *WorkerTracker) WorkerStates() []interface{} {
	wt.mu.RLock()
	states := make([]interface{}, 0, len(wt.workers))
	for _, worker := range wt.workers {
		states = append(states, worker)
	}
	wt.mu.RUnlock()
	return states
}
