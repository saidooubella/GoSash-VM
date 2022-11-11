package bytes

func Append1Byte(bytes []uint8, offset int, value uint8) {
	bytes[offset] = uint8(value)
}

func Append2Byte(bytes []uint8, offset int, value uint16) {
	bytes[offset] = uint8((value >> 8) & 0xFF)
	bytes[offset+1] = uint8(value & 0xFF)
}

func Append4Byte(bytes []uint8, offset int, value uint32) {
	bytes[offset] = uint8((value >> 24) & 0xFF)
	bytes[offset+1] = uint8((value >> 16) & 0xFF)
	bytes[offset+2] = uint8((value >> 8) & 0xFF)
	bytes[offset+3] = uint8(value & 0xFF)
}

func Append8Byte(bytes []uint8, offset int, value uint64) {
	bytes[offset] = uint8((value >> 56) & 0xFF)
	bytes[offset+1] = uint8((value >> 48) & 0xFF)
	bytes[offset+2] = uint8((value >> 40) & 0xFF)
	bytes[offset+3] = uint8((value >> 32) & 0xFF)
	bytes[offset+4] = uint8((value >> 24) & 0xFF)
	bytes[offset+5] = uint8((value >> 16) & 0xFF)
	bytes[offset+6] = uint8((value >> 8) & 0xFF)
	bytes[offset+7] = uint8(value & 0xFF)
}

func Get1Byte(bytes []uint8, offset int) uint8 {
	return bytes[offset]
}

func Get2Byte(bytes []uint8, offset int) uint16 {
	return (uint16(bytes[offset])&0xFF)<<8 |
		(uint16(bytes[offset+1]) & 0xFF)
}

func Get4Byte(bytes []uint8, offset int) uint32 {
	return (uint32(bytes[offset])&0xFF)<<24 |
		(uint32(bytes[offset+1])&0xFF)<<16 |
		(uint32(bytes[offset+2])&0xFF)<<8 |
		(uint32(bytes[offset+3]) & 0xFF)
}

func Get8Byte(bytes []uint8, offset int) uint64 {
	return (uint64(bytes[offset])&0xFF)<<56 |
		(uint64(bytes[offset+1])&0xFF)<<48 |
		(uint64(bytes[offset+2])&0xFF)<<40 |
		(uint64(bytes[offset+3])&0xFF)<<32 |
		(uint64(bytes[offset+4])&0xFF)<<24 |
		(uint64(bytes[offset+5])&0xFF)<<16 |
		(uint64(bytes[offset+6])&0xFF)<<8 |
		(uint64(bytes[offset+7]) & 0xFF)
}
