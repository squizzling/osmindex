package pb

func DecodeBytes(buf []byte, next *int) []byte {
	return ReadBytes(buf, next)
}

func ReadBytes(buf []byte, next *int) []byte {
	byteLength := int(DecodeVarInt(buf, next))
	outputBuffer := buf[*next : *next+byteLength]
	*next += byteLength
	return outputBuffer
}

func EncodeBuffer(bufOutput []byte, bufInput []byte) []byte {
	if len(bufOutput)+len(bufInput) > cap(bufOutput) {
		//panic(fmt.Sprintf("len=%d+%d, cap=%d", len(bufInput), len(bufOutput), cap(bufOutput)))
	}

	bufOutput = EncodeVarInt(bufOutput, uint64(len(bufInput)))
	bufOutput = append(bufOutput, bufInput...)
	return bufOutput
}

func SkipBytes(buf []byte, next *int) {
	byteLength := int(DecodeVarInt(buf, next))
	*next += byteLength
}

func EncodeIdBuffer(bufOutput []byte, id uint64, bufInput []byte) []byte {
	bufOutput = EncodeVarInt(bufOutput, MakeIdType(id, PbFixedBytes))
	bufOutput = EncodeBuffer(bufOutput, bufInput)
	return bufOutput
}

func DecodeStringOpt(buf []byte, next *int) *string {
	s := DecodeString(buf, next)
	return &s
}

func DecodeString(buf []byte, next *int) string {
	return string(DecodeBytes(buf, next))
}

func EncodeString(bufOutput []byte, s string) []byte {
	return EncodeBuffer(bufOutput, []byte(s))
}
