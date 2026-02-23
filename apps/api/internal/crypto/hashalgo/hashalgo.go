package hashalgo

const (
	Invalid = uint16(0)
	SHA256  = uint16(2)
	SHA384  = uint16(3)
	SHA512  = uint16(4)
)

func ToString(value uint16) string {
	switch value {
	case SHA256:
		return "SHA256"
	case SHA384:
		return "SHA384"
	case SHA512:
		return "SHA512"
	default:
		return "Invalid"
	}
}

func ToValue(value string) uint16 {
	switch value {
	case "SHA256":
		return SHA256
	case "SHA384":
		return SHA384
	case "SHA512":
		return SHA512
	default:
		return Invalid
	}
}
