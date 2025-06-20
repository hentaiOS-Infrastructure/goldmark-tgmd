package tgmd

import "unsafe"

// StringToBytes convert a string to a byte slice without copying.
//
// Note: The returned byte slice shares the same underlying data as the string.
// Modifying the slice can lead to undefined behavior.
func StringToBytes(v string) []byte {
	return unsafe.Slice(unsafe.StringData(v), len(v))
}
