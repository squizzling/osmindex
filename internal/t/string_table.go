package t

import (
	"fmt"

	"github.com/squizzling/osmindex/internal/pb"
)

type StringTable struct {
	S []string // 1
}

const (
	stringTableS = 1
)

func ReadStringTable(buf []byte, st *StringTable) {
	for next := 0; next < len(buf); {
		id := pb.DecodeVarInt(buf, &next)
		switch id {
		case pb.MakeIdType(stringTableS, pb.PbFixedBytes): // S
			st.S = append(st.S, pb.DecodeString(buf, &next))
		default:
			panic(fmt.Sprintf("StringTable: Unknown: %d (id=%d / t=%d)", id, id>>3, id&7))
		}
	}
}

func (st *StringTable) Write(buf []byte) []byte {
	for _, s := range st.S {
		buf = pb.EncodeIdBuffer(buf, stringTableS, []byte(s))
	}
	return buf
}
