package nodeindex

const (
	CountBits = 24
	CountMax = 1 << CountBits
	CountMask = CountMax - 1
	OffsetBits = 64 - CountBits
	OffsetMax = 1 << OffsetBits
)

type CountOffset uint64

func (co CountOffset) Offset() uint64 {
	return uint64(co) >> CountBits
}

func (co CountOffset) Count() uint64 {
	return uint64(co) & CountMask
}

func NewCountOffset(count, offset uint64) CountOffset {
	if offset >= OffsetMax {
		panic("offset out of range")
	}
	if count > CountMax {
		panic("too many items")
	}
	return CountOffset((offset << CountBits) | count)
}
