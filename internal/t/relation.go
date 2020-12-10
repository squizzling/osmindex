package t

import (
	"fmt"

	"github.com/squizzling/osmindex/internal/pb"
)

const (
	RelationMemberTypeNode     = 0
	RelationMemberTypeWay      = 1
	RelationMemberTypeRelation = 2
)

const (
	relationId         = 1
	relationKeys       = 2
	relationVals       = 3
	relationInfo       = 4
	relationRolesSID   = 8
	relationMemIDs     = 9
	relationMemberType = 10
)

type Relation struct {
	ID         RelId    // 1
	Keys       []uint32 // 2, packed
	Vals       []uint32 // 3, packed
	RolesSID   []int32  // 8, packed
	MemIDs     []int64  // 9, sint64, packed, delta
	MemberType []uint64 // 10, packed
}

func ReadRelation(buf []byte, r *Relation) {
	for next := 0; next < len(buf); {
		id := pb.DecodeVarInt(buf, &next)
		switch id {
		case pb.MakeIdType(relationId, pb.PbVarInt):
			r.ID = RelId(pb.DecodeI64(buf, &next))
		case pb.MakeIdType(relationKeys, pb.PbFixedBytes):
			r.Keys = pb.DecodeU32Packed(buf, &next)
		case pb.MakeIdType(relationVals, pb.PbFixedBytes):
			r.Vals = pb.DecodeU32Packed(buf, &next)
		case pb.MakeIdType(relationInfo, pb.PbFixedBytes):
			pb.SkipBytes(buf, &next)
		case pb.MakeIdType(relationRolesSID, pb.PbFixedBytes):
			r.RolesSID = pb.DecodeI32Packed(buf, &next)
		case pb.MakeIdType(relationMemIDs, pb.PbFixedBytes):
			r.MemIDs = pb.DecodeS64PackedDelta(buf, &next)
		case pb.MakeIdType(relationMemberType, pb.PbFixedBytes):
			r.MemberType = pb.DecodeU64Packed(buf, &next)
		default:
			panic(fmt.Sprintf("ReadRelation: Unknown: %d (id=%d / t=%d)", id, id>>3, id&7))
		}
	}
}

func (r *Relation) Write(buf []byte) []byte {
	buf = pb.EncodeIdVarInt(buf, relationId, uint64(r.ID))
	buf = pb.Embed16384(relationKeys, pb.EncodeU32PackedFunc(r.Keys), buf)
	buf = pb.Embed16384(relationVals, pb.EncodeU32PackedFunc(r.Vals), buf)
	buf = pb.Embed2097152(relationRolesSID, pb.EncodeI32PackedFunc(r.RolesSID), buf)
	buf = pb.Embed2097152(relationMemIDs, pb.EncodeS64PackedDeltaFunc(r.MemIDs), buf)
	buf = pb.Embed2097152(relationMemberType, pb.EncodeU64Packed(r.MemberType), buf)
	return buf
}
