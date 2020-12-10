package morton

var mortonTable256 []uint64

func init() {
	mortonTable256 = make([]uint64, 256)
	for i := uint64(0); i < 256; i++ {
		o := uint64(0)
		o |= (i & 1) << 0
		o |= (i & 2) << 1
		o |= (i & 4) << 2
		o |= (i & 8) << 3
		o |= (i & 16) << 4
		o |= (i & 32) << 5
		o |= (i & 64) << 6
		o |= (i & 128) << 7
		mortonTable256[i] = o
	}
}

func interleaveMorton16(x, y int32) uint64 {
	lo := ((mortonTable256[y&0xff]) << 1) | (mortonTable256[x&0xff])
	hi := ((mortonTable256[y>>8]) << 1) | (mortonTable256[x>>8])
	return lo | hi<<16
}

func Encode(evens, odds int32) uint64 {
	lo := interleaveMorton16(evens&0xffff, odds&0xffff)
	hi := interleaveMorton16((evens>>16)&0xffff, (odds>>16)&0xffff)
	return lo | hi<<32
}

func Decode(x uint64) (evens uint32, odds uint32) {
	e := uint32(0)
	o := uint32(0)
	for i := 0; i < 32; i++ {
		e >>= 1
		o >>= 1
		if x&1 != 0 {
			e |= 0x80000000
		}
		if x&2 != 0 {
			o |= 0x80000000
		}
		x >>= 2
	}
	return e, o
}

func DecodeLocation(x uint64) (evens int32, odds int32) {
	e, o := Decode(x)

	lon := int32(e)
	lat := int32(o)

	lat <<= 1
	lat >>= 1

	return lon, lat
}
