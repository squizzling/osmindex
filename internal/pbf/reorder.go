package pbf

import (
	"sync"
)

type reorderBuffer map[uint64]Block

func (rb reorderBuffer) pop(k uint64) (Block, bool) {
	obo, ok := rb[k]
	if ok {
		delete(rb, k)
	}
	return obo, ok
}

func reorder(inBlocks <-chan Block, outBlocks chan<- Block, wg *sync.WaitGroup) {
	rob := make(reorderBuffer)
	next := uint64(0)
	for b := range inBlocks {
		rob[b.Index] = b
		for b, ok := rob.pop(next); ok; b, ok = rob.pop(next) {
			next++
			outBlocks <- b
		}
	}
	if len(rob) != 0 {
		panic("channel closed without all blocks processed")
	}
	wg.Done()
}


