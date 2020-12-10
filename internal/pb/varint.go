package pb

func EncodeIdVarInt(buf []byte, id uint64, v uint64) []byte {
	buf = EncodeVarInt(buf, MakeIdType(id, PbVarInt))
	buf = EncodeVarInt(buf, v)
	return buf
}

func EncodeVarInt(b []byte, v uint64) []byte {
	//if len(b) + 16 > cap(b) {
	//	panic(fmt.Sprintf("len=%d, cap=%d", len(b), cap(b)))
	//}
	switch {
	case v < 1<<7:
		b = append(b, byte(v))
	case v < 1<<14:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte(v>>7))
	case v < 1<<21:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte(v>>14))
	case v < 1<<28:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte(v>>21))
	case v < 1<<35:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte(v>>28))
	case v < 1<<42:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte(v>>35))
	case v < 1<<49:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte(v>>42))
	case v < 1<<56:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte(v>>49))
	case v < 1<<63:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte(v>>56))
	default:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte((v>>56)&0x7f|0x80),
			1)
	}
	return b
}

func DecodeVarInt(b []byte, s *int) (v uint64) {
	var y uint64
	start := *s
	v = uint64(b[start])
	if v < 0x80 {
		*s += 1
		return v
	}
	v -= 0x80

	y = uint64(b[start+1])
	v += y << 7
	if y < 0x80 {
		*s += 2
		return v
	}
	v -= 0x80 << 7

	y = uint64(b[start+2])
	v += y << 14
	if y < 0x80 {
		*s += 3
		return v
	}
	v -= 0x80 << 14

	y = uint64(b[start+3])
	v += y << 21
	if y < 0x80 {
		*s += 4
		return v
	}
	v -= 0x80 << 21

	y = uint64(b[start+4])
	v += y << 28
	if y < 0x80 {
		*s += 5
		return v
	}
	v -= 0x80 << 28

	y = uint64(b[start+5])
	v += y << 35
	if y < 0x80 {
		*s += 6
		return v
	}
	v -= 0x80 << 35

	y = uint64(b[start+6])
	v += y << 42
	if y < 0x80 {
		*s += 7
		return v
	}
	v -= 0x80 << 42

	y = uint64(b[start+7])
	v += y << 49
	if y < 0x80 {
		*s += 8
		return v
	}
	v -= 0x80 << 49

	y = uint64(b[start+8])
	v += y << 56
	if y < 0x80 {
		*s += 9
		return v
	}
	v -= 0x80 << 56

	y = uint64(b[start+9])
	v += y << 63
	if y < 2 {
		*s += 10
		return v
	}
	panic("overflow")
}

func DecodeZigZag(x uint64) int64 {
	return int64(x>>1) ^ int64(x)<<63>>63
}

func EncodeZigZag(x int64) uint64 {
	return uint64(x<<1) ^ uint64(x>>63)
}
