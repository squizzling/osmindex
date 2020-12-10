package pb

// {encode,decode}{I,S,U}{32,64}

// DecodeS64 reads a single 64bit zigzag varint
func DecodeS64(buf []byte, next *int) int64 {
	return DecodeZigZag(DecodeVarInt(buf, next))
}

// DecodeI64 reads a single 64bit signed varint
func DecodeI64(buf []byte, next *int) int64 {
	return int64(DecodeVarInt(buf, next))
}

// DecodeU64 reads a single 64bit unsigned varint
func DecodeU64(buf []byte, next *int) uint64 {
	return DecodeVarInt(buf, next)
}

func DecodeS64PackedDelta(buf []byte, next *int) []int64 {
	innerBuf := DecodeBytes(buf, next)
	output := make([]int64, len(innerBuf)) // more memory but no reallocation
	idx := 0
	last := int64(0)
	for innerNext := 0; innerNext < len(innerBuf); {
		last += DecodeZigZag(DecodeVarInt(innerBuf, &innerNext))
		output[idx] = last
		idx++
	}
	return output[:idx]
}

// DecodeS64PackedDeltaZero is a variation on DecodeS64PackedDelta which decodes
// until it hits a 0, rather than using a length.
func DecodeS64PackedDeltaZero(buf []byte) []int64 {
	var output []int64
	last := int64(0)
	for next := 0; next < len(buf); {
		nextDelta := DecodeVarInt(buf, &next)
		if nextDelta == 0 {
			return output
		}
		last += DecodeZigZag(nextDelta)
		output = append(output, last)
	}
	panic("ran out of data")
}

// EncodeS64PackedDeltaZero is a variation on EncodeS64PackedDelta, which encodes
// without a length, and puts a 0 on the end.
func EncodeS64PackedDeltaZero(buf []byte, values []int64) []byte {
	last := int64(0)
	for _, value := range values {
		buf = EncodeVarInt(buf, EncodeZigZag(value-last))
		last = value
	}
	buf = append(buf, 0)
	return buf
}

func EncodeS64PackedDelta(buf []byte, values []int64) []byte {
	last := int64(0)
	var innerBuf []byte
	for _, value := range values {
		innerBuf = EncodeVarInt(innerBuf, EncodeZigZag(value-last))
		last = value
	}
	buf = EncodeBuffer(buf, innerBuf)
	return buf
}

func EncodeS64PackedDeltaFunc(values []int64) func(buf []byte) []byte {
	return func(buf []byte) []byte {
		last := int64(0)
		for _, value := range values {
			buf = EncodeVarInt(buf, EncodeZigZag(value-last))
			last = value
		}
		return buf
	}
}

func DecodeS64Packed(buf []byte, next *int) []int64 {
	innerBuf := DecodeBytes(buf, next)
	output := make([]int64, len(innerBuf)) // more memory but no reallocation
	idx := 0
	for innerNext := 0; innerNext < len(innerBuf); {
		output[idx] = DecodeZigZag(DecodeVarInt(innerBuf, &innerNext))
		idx++
	}
	return output[:idx]
}

func DecodeU64Packed(buf []byte, next *int) []uint64 {
	var output []uint64
	innerBuf := DecodeBytes(buf, next)
	for innerNext := 0; innerNext < len(innerBuf); {
		output = append(output, DecodeVarInt(innerBuf, &innerNext))
	}
	return output
}

func EncodeU64Packed(values []uint64) func(buf []byte) []byte {
	return func(buf []byte) []byte {
		for _, value := range values {
			buf = EncodeVarInt(buf, value)
		}
		return buf
	}
}

func DecodeS64Opt(buf []byte, next *int) *int64 {
	o := DecodeS64(buf, next)
	return &o
}

func DecodeI64Opt(buf []byte, next *int) *int64 {
	o := DecodeI64(buf, next)
	return &o
}

func DecodeU64Opt(buf []byte, next *int) *uint64 {
	o := DecodeU64(buf, next)
	return &o
}
