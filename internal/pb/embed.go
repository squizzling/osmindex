package pb

type Embedder func(buf []byte) []byte

func Embed128(id uint64, e Embedder, buf []byte) []byte {
	const sizeOfReservation = 1 // remember to update this
	buf = EncodeVarInt(buf, MakeIdType(id, PbFixedBytes))
	startIndex := len(buf)
	buf = Reserve128(buf)
	buf = e(buf)
	EmbedInt128(buf[startIndex:], len(buf)-(startIndex+sizeOfReservation))
	return buf
}

func Embed16384(id uint64, e Embedder, buf []byte) []byte {
	const sizeOfReservation = 2 // remember to update this
	buf = EncodeVarInt(buf, MakeIdType(id, PbFixedBytes))
	startIndex := len(buf)
	buf = Reserve16384(buf)
	buf = e(buf)
	EmbedInt16384(buf[startIndex:], len(buf)-(startIndex+sizeOfReservation))
	return buf
}

func Embed2097152(id uint64, e Embedder, buf []byte) []byte {
	const sizeOfReservation = 3 // remember to update this
	buf = EncodeVarInt(buf, MakeIdType(id, PbFixedBytes))
	startIndex := len(buf)
	buf = Reserve2097152(buf)
	buf = e(buf)
	EmbedInt2097152(buf[startIndex:], len(buf)-(startIndex+sizeOfReservation))
	return buf
}

func Embed268435456(id uint64, e Embedder, buf []byte) []byte {
	const sizeOfReservation = 4 // remember to update this
	buf = EncodeVarInt(buf, MakeIdType(id, PbFixedBytes))
	startIndex := len(buf)
	buf = Reserve268435456(buf)
	buf = e(buf)
	EmbedInt268435456(buf[startIndex:], len(buf)-(startIndex+sizeOfReservation))
	return buf
}

func Reserve128(b []byte) []byte {
	return append(b, 0)
}

func Reserve16384(b []byte) []byte {
	return append(b, 0, 0)
}

func Reserve2097152(b []byte) []byte {
	return append(b, 0, 0, 0)
}

func Reserve268435456(b []byte) []byte {
	return append(b, 0, 0, 0, 0)
}

func EmbedInt128(b []byte, v int) {
	if v >= 128 {
		panic("128 overflow")
	}
	b[0] = byte(v)
}

func EmbedInt16384(b []byte, v int) {
	if v >= 16384 {
		panic("16384 overflow")
	}
	b[0] = byte((v>>0)&0x7f | 0x80)
	b[1] = byte(v >> 7)
}

func EmbedInt2097152(b []byte, v int) {
	if v >= 2097152 {
		panic("2097152 overflow")
	}
	b[0] = byte((v>>0)&0x7f | 0x80)
	b[1] = byte((v>>7)&0x7f | 0x80)
	b[2] = byte(v >> 14)
}

func EmbedInt268435456(b []byte, v int) {
	if v >= 268435456 {
		panic("268435456 overflow")
	}
	b[0] = byte((v>>0)&0x7f | 0x80)
	b[1] = byte((v>>7)&0x7f | 0x80)
	b[2] = byte((v>>14)&0x7f | 0x80)
	b[3] = byte(v >> 21)
}
