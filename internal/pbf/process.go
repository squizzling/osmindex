package pbf

import (
	"context"
	"encoding/binary"
	"sync"
	"sync/atomic"

	"github.com/squizzling/osmindex/internal/iio"
	"github.com/squizzling/osmindex/internal/t"
)

type Block struct {
	Index    uint64
	Data     interface{} // can't wait for generics...
	BlobType string
}

type WorkFunc func(ctx context.Context, inBlocks <-chan Block, outBlocks chan<- Block)
type WriteFunc func(ctx context.Context, inBlocks <-chan Block)

type ProcessFileState struct {
	currentState uint64
	at           uint64
	max          uint64
	done         chan struct{}
	filename     string
}

const (
	ProcessFileStateInitializing = iota
	ProcessFileStateLoadingBlock
	ProcessFileStateSendingBlock
	ProcessFileStateShuttingDown
)

func (pfs *ProcessFileState) SetState(e uint64) {
	atomic.StoreUint64(&pfs.currentState, e)
}

func (pfs *ProcessFileState) CurrentState() uint64 {
	return atomic.LoadUint64(&pfs.currentState)
}

func (pfs *ProcessFileState) At() uint64 {
	return atomic.LoadUint64(&pfs.at)
}

func (pfs *ProcessFileState) Max() uint64 {
	return pfs.max
}

func (pfs *ProcessFileState) Wait() {
	<-pfs.done
}

func (pfs *ProcessFileState) Filename() string {
	return pfs.filename
}

func (pfs *ProcessFileState) processFileWork(
	ctx context.Context,
	file string,
	workStep WorkFunc,
	writeStep WriteFunc,
	workers int,
) {
	inputFile := iio.OsOpen(file)
	inputBuffer := iio.MMapMap(inputFile)

	var wgWriter sync.WaitGroup
	var wgReorder sync.WaitGroup

	var chReorderBlobs chan Block
	var chOutputBlobs chan Block

	if writeStep != nil {
		chOutputBlobs = make(chan Block)
		wgWriter.Add(1)
		go func() {
			writeStep(ctx, chOutputBlobs)
			wgWriter.Done()
		}()

		wgReorder.Add(1)
		chReorderBlobs = make(chan Block)
		go reorder(chReorderBlobs, chOutputBlobs, &wgReorder)
	}

	var wgWorker sync.WaitGroup
	chInputBlobs := make(chan Block)
	wgWorker.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			workStep(ctx, chInputBlobs, chReorderBlobs)
			wgWorker.Done()
		}()
	}

	atomic.StoreUint64(&pfs.at, 0)
	pfs.max = uint64(len(inputBuffer))

	blobIndex := uint64(0)

	var pbr t.PBReader

	for pfs.at < pfs.max {
		pfs.SetState(ProcessFileStateLoadingBlock)
		blobHeaderSize := uint64(binary.BigEndian.Uint32(inputBuffer[pfs.at:]))
		atomic.AddUint64(&pfs.at, 4)

		var bh t.BlobHeader
		pbr.ReadBlobHeader(inputBuffer[pfs.at:pfs.at+blobHeaderSize], &bh)
		atomic.AddUint64(&pfs.at, blobHeaderSize)

		if bh.Type != "OSMHeader" && bh.Type != "OSMData" {
			panic(bh.Type)
		}

		obh := Block{
			Index:    blobIndex,
			Data:     inputBuffer[pfs.at : pfs.at+uint64(bh.DataSize)],
			BlobType: bh.Type,
		}

		pfs.SetState(ProcessFileStateSendingBlock)

		select {
		case <-ctx.Done():
			goto cleanup
		case chInputBlobs <- obh:
		}

		blobIndex++
		atomic.AddUint64(&pfs.at, uint64(bh.DataSize))
	}

cleanup:
	pfs.SetState(ProcessFileStateShuttingDown)
	close(chInputBlobs)
	wgWorker.Wait()
	if chReorderBlobs != nil {
		close(chReorderBlobs)
		wgReorder.Wait()
	}
	if chOutputBlobs != nil {
		close(chOutputBlobs)
		wgWriter.Wait()
	}
	iio.MMapUnmap(inputBuffer)
	iio.FClose(inputFile)
	close(pfs.done)
}

func ProcessFileAsync(
	ctx context.Context,
	file string,
	workStep WorkFunc,
	writeStep WriteFunc,
	workers int,
) *ProcessFileState {
	ps := &ProcessFileState{
		filename:     file,
		currentState: uint64(ProcessFileStateInitializing),
		done:         make(chan struct{}),
	}
	go ps.processFileWork(
		ctx,
		file,
		workStep,
		writeStep,
		workers,
	)
	return ps
}
