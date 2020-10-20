package util

import "encoding/binary"

func SliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ToByteArray(value int64) []byte {
	res := make([]byte, 8)
	binary.LittleEndian.PutUint64(res, uint64(value))
	return res
}

func ToInt64(value []byte) int64 {
	return int64(binary.LittleEndian.Uint64(value))
}
