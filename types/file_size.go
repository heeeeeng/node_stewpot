package types

//type FileSize int64

//func (fs FileSize) ToInt64() int64 {
//	return int64(fs)
//}

const (
	Bit  = int64(1)
	Byte = 8 * Bit
	KB   = 1024 * Byte
	MB   = 1024 * KB
	GB   = 1024 * MB
)
