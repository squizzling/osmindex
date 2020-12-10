package wayindex

import (
	"fmt"
)

type WayIndexData uint64

func (id WayIndexData) WayID() int64 {
	// Will not be negative
	return int64(id >> 33)
}

func (id WayIndexData) Offset() int64 {
	return int64(id) & 0x1_ffff_ffff
}

func MakeIndexData(wayID int64, offset int64) WayIndexData {
	if wayID > (1 << 31) {
		panic("wayID too large")
	}
	if offset > (1 << 33) {
		panic(fmt.Sprintf("offset too large: %d\n", offset))
	}
	return WayIndexData(wayID << 33 | offset)
}
