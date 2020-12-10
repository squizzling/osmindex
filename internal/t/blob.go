package t

import (
	"bytes"
	"compress/zlib"
	"fmt"

	"github.com/squizzling/osmindex/internal/pb"
	"github.com/squizzling/osmindex/internal/pool"
)

const (
	blobRaw      = 1
	blobRawSize  = 2
	blobZLibData = 3
)

type Blob struct {
	Raw      []byte // 1, optional
	RawSize  *int32 // 2, optional
	ZLibData []byte // 3, optional

	RawFunc pb.Embedder
}

// TODO: Assess making this return a Blob and remove *Blob from the code
func (pbr *PBReader) ReadBlob(buf []byte, b *Blob) {
	for next := 0; next < len(buf); {
		id := pb.DecodeVarInt(buf, &next)
		switch id {
		case pb.MakeIdType(blobRaw, pb.PbFixedBytes): // Raw
			b.Raw = pb.DecodeBytes(buf, &next)
		case pb.MakeIdType(blobRawSize, pb.PbVarInt): // RawSize
			b.RawSize = pb.DecodeI32Opt(buf, &next)
		case pb.MakeIdType(blobZLibData, pb.PbFixedBytes): // ZLibData
			b.ZLibData = pb.DecodeBytes(buf, &next)
		default:
			panic(fmt.Sprintf("Blob: Unknown: %d (id=%d / t=%d)", id, id>>3, id&7))
		}
	}
}

func (b *Blob) Write(buf []byte) []byte {
	// If RawFunc is provided, invoke it, otherwise use Raw
	if b.RawFunc != nil {
		buf = pb.Embed268435456(blobRaw, b.RawFunc, buf)
	} else if len(b.Raw) > 0 {
		buf = pb.EncodeIdBuffer(buf, blobRaw, b.Raw)
	}
	if b.RawSize != nil {
		buf = pb.EncodeIdVarInt(buf, blobRawSize, uint64(*b.RawSize))
	}
	if len(b.ZLibData) > 0 {
		buf = pb.EncodeIdBuffer(buf, blobZLibData, b.ZLibData)
	}
	return buf
}

func (b *Blob) GetRawData() pool.BorrowedBuffer {
	if b.Raw != nil {
		return pool.NonBorrowedBuffer(b.Raw)
	} else if b.ZLibData != nil {
		return b.decodeZLib()
	} else {
		panic("unexpected compression")
	}
}

func (b *Blob) decodeZLib() pool.BorrowedBuffer {
	buf := bytes.NewBuffer(b.ZLibData)
	zlibReader, err := zlib.NewReader(buf)
	if err != nil {
		panic(err)
	}

	borrowedBuffer := pool.ByteSlice.Get(uint64(*b.RawSize))
	outBuf := borrowedBuffer.Buffer

	remaining := int(*b.RawSize)
	start := 0
	for remaining > 0 {
		remainingBuffer := outBuf[start:cap(outBuf)]
		read, err := zlibReader.Read(remainingBuffer)
		remaining -= read
		start += read
		if remaining == 0 {
			return borrowedBuffer
		}
		if err != nil {
			panic(err)
		}
	}

	return borrowedBuffer
}
