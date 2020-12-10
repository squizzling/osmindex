package pb

type pbType uint64

const (
	PbVarInt     = pbType(0)
	PbFixedBytes = pbType(2)
)

func MakeIdType(id uint64, t pbType) uint64 {
	return id<<3 | uint64(t)
}
