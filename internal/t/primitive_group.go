package t

import (
	"fmt"

	"github.com/squizzling/osmindex/internal/pb"
)

const (
	primitiveGroupNodes      = 1
	primitiveGroupDenseNodes = 2
	primitiveGroupWays       = 3
	primitiveGroupRelations  = 4
	primitiveGroupChangeSets = 5
)

type PrimitiveGroup struct {
	//Nodes        []*Node // 1

	RawDense     []byte
	Dense        *DenseNodes // 2
	Ways         []*Way      // 3
	RawWays      [][]byte
	Relations    []*Relation // 4
	RawRelations [][]byte

	//ChangeSets   []*ChangeSet // 5
}

func (pbr *PBReader) ReadPrimitiveGroup(buf []byte, pg *PrimitiveGroup) {
	for next := 0; next < len(buf); {
		id := pb.DecodeVarInt(buf, &next)
		switch id {
		case pb.MakeIdType(primitiveGroupNodes, pb.PbFixedBytes):
			panic("non-dense nodes not supported")
		case pb.MakeIdType(primitiveGroupDenseNodes, pb.PbFixedBytes):
			if pbr.SkipDenseNodes {
				pg.RawDense = pb.DecodeBytes(buf, &next)
			} else {
				pg.Dense = &DenseNodes{}
				ReadDenseNodes(pb.DecodeBytes(buf, &next), pg.Dense)
			}
		case pb.MakeIdType(primitiveGroupWays, pb.PbFixedBytes):
			if pbr.SkipWays {
				pg.RawWays = append(pg.RawWays, pb.DecodeBytes(buf, &next))
			} else {
				w := &Way{}
				pbr.ReadWay(pb.DecodeBytes(buf, &next), w)
				pg.Ways = append(pg.Ways, w)
			}
		case pb.MakeIdType(primitiveGroupRelations, pb.PbFixedBytes):
			if pbr.SkipRelations {
				pg.RawRelations = append(pg.RawRelations, pb.DecodeBytes(buf, &next))
			} else {
				r := &Relation{}
				ReadRelation(pb.DecodeBytes(buf, &next), r)
				pg.Relations = append(pg.Relations, r)
			}
		case pb.MakeIdType(primitiveGroupChangeSets, pb.PbFixedBytes):
			pb.SkipBytes(buf, &next) // I assume this is sufficient, I have nothing to test on.
		default:
			panic(fmt.Sprintf("PrimitiveGroup: Unknown: %d (id=%d / t=%d)", id, id>>3, id&7))
		}
	}
}

func (pg *PrimitiveGroup) Write(buf []byte) []byte {
	if pg.Dense != nil {
		buf = pb.Embed268435456(primitiveGroupDenseNodes, pg.Dense.Write, buf)
	} else if pg.RawDense != nil {
		buf = pb.EncodeIdBuffer(buf, primitiveGroupDenseNodes, pg.RawDense)
	}

	if pg.Ways != nil {
		for _, w := range pg.Ways {
			buf = pb.Embed16384(primitiveGroupWays, w.Write, buf)
		}
	} else {
		for _, w := range pg.RawWays {
			buf = pb.EncodeIdBuffer(buf, primitiveGroupWays, w)
		}
	}

	if pg.Relations != nil {
		for _, r := range pg.Relations {
			buf = pb.Embed2097152(primitiveGroupRelations, r.Write, buf)
		}
	} else {
		for _, r := range pg.RawRelations {
			buf = pb.EncodeIdBuffer(buf, primitiveGroupRelations, r)
		}
	}

	return buf
}
