package pool

import (
	"sync"
	"sync/atomic"
)

type ByteSlicePool struct {
	p       sync.Pool
	defSize uint64
}

var ByteSlice = NewByteSlicePool(65536)

func NewByteSlicePool(bufferSize uint64) *ByteSlicePool {
	bsp := &ByteSlicePool{
		defSize: bufferSize,
	}
	bsp.p = sync.Pool{
		New: func() interface{} {
			return make([]byte, atomic.LoadUint64(&bsp.defSize))
		},
	}
	return bsp
}

func (bsp *ByteSlicePool) Get(sz uint64) BorrowedBuffer {
	if sz > atomic.LoadUint64(&bsp.defSize) {
		// doesn't need to be a true CAS loop, this is good enough
		atomic.StoreUint64(&bsp.defSize, sz)
	}

	buf := bsp.p.Get().([]byte)
	if cap(buf) < int(sz) {
		// discard the thing we pulled from the pool, because we don't want to keep cycling it through
		buf = make([]byte, sz)
	}

	return BorrowedBuffer{
		Buffer:     buf[:sz],
		isBorrowed: true,
	}
}

func (bsp *ByteSlicePool) Put(bb BorrowedBuffer) {
	if bb.isBorrowed {
		bsp.p.Put(bb.Buffer)
	}
}

func NonBorrowedBuffer(buf []byte) BorrowedBuffer {
	return BorrowedBuffer{
		Buffer:     buf,
		isBorrowed: false,
	}
}

type BorrowedBuffer struct {
	Buffer     []byte
	isBorrowed bool
}
