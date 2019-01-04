package mcpegql

import "encoding/hex"

// ConvertKey takes a byte array and returns a string if all characters are printable (else "")  hex-string-encoded versions of key
func ConvertKey(k []byte) (stringKey, hexKey string) {
	allAscii := true
	for i := range k {
		if !isAscii(k[i]) {
			allAscii = false
		}
	}
	if allAscii {
		stringKey = string(k[:])
	}
	hexKey = hex.EncodeToString(k)
	return
}

func isAscii(b byte) bool {
	return b >= 0x20 && b <= 0x7e
}

// If key is certain length and x/z MSBs aren't both printable ASCII, assume chunk key (not ideal, but probably works in all real cases)
func IsChunkKey(k []byte) bool {
	isChunk := false
	for _, e := range []int{9, 10, 13, 14} {
		if e == len(k) {
			if !(isAscii(k[3]) && isAscii(k[7])) {
				isChunk = true
			}
		}
	}
	return isChunk
}
