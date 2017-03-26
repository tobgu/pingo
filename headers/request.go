package headers

type RequestVersion [8]byte

var v1 = RequestVersion{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x01}

type RequestV1 struct {
	RequestBodyLen  uint32
	ResponseBodyLen uint32
	ConfigBitmap    uint32
}
