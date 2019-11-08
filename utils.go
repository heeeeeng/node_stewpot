package main

func BytesToInt64(b []byte) int64 {
	return int64(b[7]) |
		int64(b[6])<<8 |
		int64(b[5])<<16 |
		int64(b[4])<<24 |
		int64(b[3])<<32 |
		int64(b[2])<<40 |
		int64(b[1])<<48 |
		int64(b[0])<<56
}

func Int64ToBytes(i int64) []byte {
	b := make([]byte, 8)

	b[7] = byte(i)
	b[6] = byte(i >> 8)
	b[5] = byte(i >> 16)
	b[4] = byte(i >> 24)
	b[3] = byte(i >> 32)
	b[2] = byte(i >> 40)
	b[1] = byte(i >> 48)
	b[0] = byte(i >> 56)

	return b
}
