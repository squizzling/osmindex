package pb

func DecodeS32(buf []byte, next *int) int32 {
	return int32(DecodeZigZag(DecodeVarInt(buf, next)))
}

func DecodeI32(buf []byte, next *int) int32 {
	return int32(DecodeVarInt(buf, next))
}

func DecodeU32(buf []byte, next *int) uint32 {
	return uint32(DecodeVarInt(buf, next))
}

func DecodeS32Packed(buf []byte, next *int) []int32 {
	innerBuf := DecodeBytes(buf, next)
	output := make([]int32, len(innerBuf))
	idx := 0
	for innerNext := 0; innerNext < len(innerBuf); {
		output[idx] = int32(DecodeZigZag(DecodeVarInt(innerBuf, &innerNext)))
		idx++
	}
	return output[:idx]
}

func EncodeS32PackedDelta(buf []byte, values []int32) []byte {
	last := int32(0)
	var innerBuf []byte
	for _, value := range values {
		innerBuf = EncodeVarInt(innerBuf, EncodeZigZag(int64(value-last)))
		last = value
	}
	buf = EncodeBuffer(buf, innerBuf)
	return buf
}

func EncodeS32Packed(buf []byte, values []int32) []byte {
	var innerBuf []byte
	for _, value := range values {
		innerBuf = EncodeVarInt(innerBuf, EncodeZigZag(int64(value)))
	}
	buf = EncodeBuffer(buf, innerBuf)
	return buf
}

func DecodeI32Packed(buf []byte, next *int) []int32 {
	innerBuf := DecodeBytes(buf, next)
	output := make([]int32, len(innerBuf))
	idx := 0
	for innerNext := 0; innerNext < len(innerBuf); {
		output[idx] = int32(DecodeVarInt(innerBuf, &innerNext))
		idx++
	}
	return output[:idx]
}

func EncodeI32Packed(buf []byte, values []int32) []byte {
	var innerBuf []byte
	for _, value := range values {
		innerBuf = EncodeVarInt(innerBuf, uint64(value))
	}
	buf = EncodeBuffer(buf, innerBuf)
	return buf
}

func EncodeI32PackedFunc(values []int32) func(buf []byte) []byte {
	return func(buf []byte) []byte {
		for _, value := range values {
			buf = EncodeVarInt(buf, uint64(value))
		}
		return buf
	}
}

func DecodeU32Packed(buf []byte, next *int) []uint32 {
	innerBuf := DecodeBytes(buf, next)
	output := make([]uint32, len(innerBuf))
	idx := 0
	for innerNext := 0; innerNext < len(innerBuf); {
		output[idx] = uint32(DecodeVarInt(innerBuf, &innerNext))
		idx++
	}
	return output[:idx]
}

func EncodeU32PackedFunc(values []uint32) func(buf []byte) []byte {
	return func(buf []byte) []byte {
		for _, value := range values {
			buf = EncodeVarInt(buf, uint64(value))
		}
		return buf
	}
}

func DecodeS32Opt(buf []byte, next *int) *int32 {
	n := DecodeS32(buf, next)
	return &n
}

func DecodeI32Opt(buf []byte, next *int) *int32 {
	n := DecodeI32(buf, next)
	return &n
}

func DecodeU32Opt(buf []byte, next *int) *uint32 {
	n := DecodeU32(buf, next)
	return &n
}
