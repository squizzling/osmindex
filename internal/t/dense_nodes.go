package t

import (
	"fmt"

	"github.com/squizzling/osmindex/internal/pb"
)

const (
	denseNodesId        = 1
	denseNodesDenseInfo = 5
	denseNodesLat       = 8
	denseNodesLon       = 9
	denseNodesKeyVals   = 10
)

type DenseNodes struct {
	Id      []int64 // 1, sint, packed, delta
	Lat     []int64 // 8, sint, packed, delta
	Lon     []int64 // 9, sint, packed, delta
	KeyVals []int32 // 10, packed

	cacheKeys [][]uint32
	cacheVals [][]uint32
}

func ReadDenseNodes(buf []byte, dn *DenseNodes) {
	for next := 0; next < len(buf); {
		id := pb.DecodeVarInt(buf, &next)
		switch id {
		case pb.MakeIdType(denseNodesId, pb.PbFixedBytes):
			dn.Id = pb.DecodeS64PackedDelta(buf, &next)
		case pb.MakeIdType(denseNodesDenseInfo, pb.PbFixedBytes):
			pb.SkipBytes(buf, &next)
		case pb.MakeIdType(denseNodesLat, pb.PbFixedBytes):
			dn.Lat = pb.DecodeS64PackedDelta(buf, &next)
		case pb.MakeIdType(denseNodesLon, pb.PbFixedBytes):
			dn.Lon = pb.DecodeS64PackedDelta(buf, &next)
		case pb.MakeIdType(denseNodesKeyVals, pb.PbFixedBytes):
			dn.KeyVals = pb.DecodeI32Packed(buf, &next)
		default:
			panic(fmt.Sprintf("ReadDenseNodes: Unknown: %d (id=%d / t=%d)", id, id>>3, id&7))
		}
	}
}

func (dn *DenseNodes) Write(buf []byte) []byte {
	buf = pb.Embed2097152(denseNodesId, pb.EncodeS64PackedDeltaFunc(dn.Id), buf)
	buf = pb.Embed2097152(denseNodesLat, pb.EncodeS64PackedDeltaFunc(dn.Lat), buf)
	buf = pb.Embed2097152(denseNodesLon, pb.EncodeS64PackedDeltaFunc(dn.Lon), buf)
	if len(dn.KeyVals) > 0 {
		buf = pb.Embed268435456(denseNodesKeyVals, pb.EncodeI32PackedFunc(dn.KeyVals), buf)
	}
	return buf
}

func (dn *DenseNodes) GetKeyVals(idx int) ([]uint32, []uint32) {
	if len(dn.KeyVals) == 0 {
		return nil, nil
	}
	if len(dn.cacheKeys) == 0 {
		startOfKeyVals := 0
		curKeys := make([]uint32, 0)
		curVals := make([]uint32, 0)
		for ; startOfKeyVals < len(dn.KeyVals); startOfKeyVals++ {
			if dn.KeyVals[startOfKeyVals] == 0 {
				dn.cacheKeys = append(dn.cacheKeys, curKeys)
				dn.cacheVals = append(dn.cacheVals, curVals)
				curKeys = make([]uint32, 0)
				curVals = make([]uint32, 0)
			} else {
				curKeys = append(curKeys, uint32(dn.KeyVals[startOfKeyVals+0]))
				curVals = append(curVals, uint32(dn.KeyVals[startOfKeyVals+1]))
				startOfKeyVals++
			}
		}
	}

	return dn.cacheKeys[idx], dn.cacheVals[idx]
}
