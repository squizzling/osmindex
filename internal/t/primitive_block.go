package t

import (
	"fmt"

	"github.com/squizzling/osmindex/internal/pb"
)

type PrimitiveBlock struct {
	StringTableRaw  []byte
	StringTable     *StringTable      // 1
	PrimitiveGroup  []*PrimitiveGroup // 2
	Granularity     int32             // 17 default 100
	DateGranularity int32             // 18 default 1000
	LatOffset       int64             // 19 default 0
	LonOffset       int64             // 20 default 0
}

const (
	primitiveBlockStringTable     = 1
	primitiveBlockPrimitiveGroup  = 2
	primitiveBlockGranularity     = 17
	primitiveBlockDateGranularity = 18
	primitiveBlockLatOffset       = 19
	primitiveBlockLonOffset       = 20
)

func (pbr *PBReader) ReadPrimitiveBlock(buf []byte, pblk *PrimitiveBlock) {
	pblk.Granularity = 100
	pblk.DateGranularity = 1000

	for next := 0; next < len(buf); {
		id := pb.DecodeVarInt(buf, &next)
		switch id {
		case pb.MakeIdType(primitiveBlockStringTable, pb.PbFixedBytes):
			if pbr.SkipStringTable {
				pblk.StringTableRaw = pb.DecodeBytes(buf, &next)
			} else {
				pblk.StringTable = &StringTable{}
				ReadStringTable(pb.DecodeBytes(buf, &next), pblk.StringTable)
			}
		case pb.MakeIdType(primitiveBlockPrimitiveGroup, pb.PbFixedBytes):
			pg := &PrimitiveGroup{}
			pbr.ReadPrimitiveGroup(pb.DecodeBytes(buf, &next), pg)
			pblk.PrimitiveGroup = append(pblk.PrimitiveGroup, pg)
		case pb.MakeIdType(primitiveBlockGranularity, pb.PbVarInt):
			pblk.Granularity = pb.DecodeI32(buf, &next)
		case pb.MakeIdType(primitiveBlockDateGranularity, pb.PbVarInt):
			pblk.DateGranularity = pb.DecodeI32(buf, &next)
		case pb.MakeIdType(primitiveBlockLatOffset, pb.PbVarInt):
			pblk.LatOffset = pb.DecodeI64(buf, &next)
		case pb.MakeIdType(primitiveBlockLonOffset, pb.PbVarInt):
			pblk.LonOffset = pb.DecodeI64(buf, &next)
		default:
			panic(fmt.Sprintf("PrimitiveBlock: Unknown: %d (id=%d / t=%d)", id, id>>3, id&7))
		}
	}
}

func (pblk *PrimitiveBlock) Write(buf []byte) []byte {
	if pblk.StringTable != nil {
		buf = pb.Embed268435456(primitiveBlockStringTable, pblk.StringTable.Write, buf)
	} else if pblk.StringTableRaw != nil {
		buf = pb.EncodeIdBuffer(buf, primitiveBlockStringTable, pblk.StringTableRaw)
	}
	for _, pg := range pblk.PrimitiveGroup {
		buf = pb.Embed268435456(primitiveBlockPrimitiveGroup, pg.Write, buf)
	}

	if pblk.Granularity != 100 {
		buf = pb.EncodeIdVarInt(buf, primitiveBlockGranularity, uint64(pblk.Granularity))
	}
	if pblk.DateGranularity != 1000 {
		buf = pb.EncodeIdVarInt(buf, primitiveBlockDateGranularity, uint64(pblk.DateGranularity))
	}
	if pblk.LatOffset != 0 {
		buf = pb.EncodeIdVarInt(buf, primitiveBlockLatOffset, uint64(pblk.LatOffset))
	}
	if pblk.LonOffset != 0 {
		buf = pb.EncodeIdVarInt(buf, primitiveBlockLonOffset, uint64(pblk.LonOffset))
	}
	return buf
}

func MakePrimitiveBlock(ns Nodes, ws Ways, rs Relations) *PrimitiveBlock {
	stb := NewStringTableBuilder()
	ns.AddTagsToStringTable(stb)
	ws.AddTagsToStringTable(stb)
	rs.AddTagsToStringTable(stb)
	stb.Finalize()

	pblk := &PrimitiveBlock{
		StringTable:     stb.ToStringTable(),
		PrimitiveGroup:  nil,
		Granularity:     100,
		DateGranularity: 1000,
		LatOffset:       0,
		LonOffset:       0,
	}

	if len(ns) > 0 {
		pblk.PrimitiveGroup = append(pblk.PrimitiveGroup, ns.ToPrimitiveGroup(stb))
	}
	if len(ws) > 0 {
		pblk.PrimitiveGroup = append(pblk.PrimitiveGroup, ws.ToPrimitiveGroup(stb))
	}
	if len(rs) > 0 {
		pblk.PrimitiveGroup = append(pblk.PrimitiveGroup, rs.ToPrimitiveBlock(stb))
	}
	return pblk
}

func (pblk *PrimitiveBlock) GetString(id uint32) string {
	return pblk.StringTable.S[id]
}
