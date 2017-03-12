package world

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

// Key is the default JSON response object
type Key struct {
	key       []byte
	keys      [][]byte
	Type      string `json:"type,omitempty"`
	StringKey string `json:"stringKey,omitempty"`
	HexKey    string `json:"hexKey,omitempty"`
	X         int    `json:"c,omitempty"`
	Z         int    `json:"z,omitempty"`
	Y         int    `json:"y,omitempty"`
}

// convertKey takes a byte array and returns a string if all characters are printable (else "")  hex-string-encoded versions of key
func convertKey(k []byte) (stringKey, hexKey string) {
	allAscii := true
	for i := range k {
		if k[i] < 0x20 || k[i] > 0x7e {
			allAscii = false
		}
	}
	if allAscii {
		stringKey = string(k[:])
	}
	hexKey = hex.EncodeToString(k)
	return
}

// TODO: Handle errors
func keyToCoords(key []byte) (x, z, y int) {
	var myx, myz int32
	var myy byte
	buf := bytes.NewReader(key[0:4])
	_ = binary.Read(buf, binary.LittleEndian, &myx)
	buf = bytes.NewReader(key[4:8])
	_ = binary.Read(buf, binary.LittleEndian, &myz)
	myy = key[len(key)-1]
	x = int(myx)
	z = int(myz)
	y = int(myy)
	return
}

func isPlayer(key []byte) bool {
	// "~local_player"
	if bytes.Equal(key, []byte{0x7e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72}[:]) {
		return true
	}
	// "player_" TODO: double-check this
	if bytes.Equal(key[:7], []byte{0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f}[:]) {
		return true
	}
	return false
}

// KeyInfo returns a struct with info deduced from just the key
func KeyInfo(key []byte) Key {
	outkey := Key{}
	outkey.StringKey, outkey.HexKey = convertKey(key)
	// "~local_player"
	if bytes.Equal(key, []byte{0x7e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72}[:]) {
		outkey.Type = "Player"
	}
	// "player_" TODO: double-check this
	if bytes.Equal(key[:7], []byte{0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f}[:]) {
		outkey.Type = "Player"
	}
	if outkey.Type == "" {
		outkey.Type = "Unknown"
	}
	return outkey
}
