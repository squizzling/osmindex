package t

import (
	"fmt"

	"github.com/squizzling/osmindex/internal/pb"
)

type BlobHeader struct {
	Type     string // 1, required
	DataSize int32  // 3, required
}

const (
	blobHeaderType     = 1
	blobHeaderDataSize = 3
)

// TODO: Assess making this return a BlobHeader and remove *BlobHeader from the code
func (pbr *PBReader) ReadBlobHeader(buf []byte, bh *BlobHeader) {
	next := 0
	for next < len(buf) {
		id := pb.DecodeVarInt(buf, &next)
		switch id {
		case pb.MakeIdType(blobHeaderType, pb.PbFixedBytes): // Type
			bh.Type = pb.DecodeString(buf, &next)
		case pb.MakeIdType(blobHeaderDataSize, pb.PbVarInt): // DataSize
			bh.DataSize = pb.DecodeI32(buf, &next)
		default:
			panic(fmt.Sprintf("BlobHeader: Unknown: %d (id=%d / t=%d)", id, id>>3, id&7))
		}
	}
}

func (bh *BlobHeader) Write(buf []byte) []byte {
	if bh.Type != "" {
		buf = pb.EncodeIdBuffer(buf, blobHeaderType, []byte(bh.Type))
	}
	if bh.DataSize != 0 {
		buf = pb.EncodeIdVarInt(buf, blobHeaderDataSize, uint64(bh.DataSize))
	}
	return buf
}
