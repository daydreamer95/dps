package common

import "encoding/binary"

func ValidateBytesSize(input []byte, maxSize int) bool {
	if binary.Size(input) > maxSize {
		return false
	}
	return true
}
