package tlv

type TLVTag uint32

const (
	Prefix_Nesting TLVTag = 0xAA << 24
	Prefix_Raw     TLVTag = 0x55 << 24
)

func NestingTag(t uint16) TLVTag {
	return TLVTag(uint32(t) | uint32(Prefix_Nesting))
}

func RawTag(t uint16) TLVTag {
	return TLVTag(uint32(t) | uint32(Prefix_Raw))
}

func IsNestingTag(t TLVTag) bool {
	return t&Prefix_Nesting == Prefix_Nesting
}

func IsRawTag(t TLVTag) bool {
	return t&Prefix_Raw == Prefix_Raw
}
