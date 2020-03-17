package types

//type FileSize int64

//func (fs FileSize) ToInt64() int64 {
//	return int64(fs)
//}

const (
	SizeBit  = int64(1)
	SizeByte = 8 * SizeBit
	SizeKB   = 1024 * SizeByte
	SizeMB   = 1024 * SizeKB
	SizeGB   = 1024 * SizeMB
)
