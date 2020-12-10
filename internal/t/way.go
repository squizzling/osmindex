package t

import (
	"fmt"

	"github.com/squizzling/osmindex/internal/pb"
)

type WayId int64

const (
	wayId   = 1
	wayKeys = 2
	wayVals = 3
	wayInfo = 4
	wayRefs = 8
)

type Way struct {
	ID   WayId    // 1
	Keys []uint32 // 2, packed
	Vals []uint32 // 3, packed
	Refs []int64  // 8, sint, packed, delta
}

func (pbr *PBReader) ReadWay(buf []byte, w *Way) {
	for next := 0; next < len(buf); {
		id := pb.DecodeVarInt(buf, &next)
		switch id {
		case pb.MakeIdType(wayId, pb.PbVarInt):
			w.ID = WayId(pb.DecodeI64(buf, &next))
		case pb.MakeIdType(wayKeys, pb.PbFixedBytes):
			w.Keys = pb.DecodeU32Packed(buf, &next)
		case pb.MakeIdType(wayVals, pb.PbFixedBytes):
			w.Vals = pb.DecodeU32Packed(buf, &next)
		case pb.MakeIdType(wayInfo, pb.PbFixedBytes):
			pb.SkipBytes(buf, &next)
		case pb.MakeIdType(wayRefs, pb.PbFixedBytes):
			w.Refs = pb.DecodeS64PackedDelta(buf, &next)
		default:
			panic(fmt.Sprintf("Way: Unknown: %d (id=%d / t=%d)", id, id>>3, id&7))
		}
	}
}

func (w *Way) Write(buf []byte) []byte {
	buf = pb.EncodeIdVarInt(buf, wayId, uint64(w.ID))
	buf = pb.Embed16384(wayKeys, pb.EncodeU32PackedFunc(w.Keys), buf)
	buf = pb.Embed16384(wayVals, pb.EncodeU32PackedFunc(w.Vals), buf)
	buf = pb.Embed16384(wayRefs, pb.EncodeS64PackedDeltaFunc(w.Refs), buf)
	return buf
}
