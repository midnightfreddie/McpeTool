package world

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
)

// Key is the default JSON response object
type Key struct {
	key       []byte
	keys      [][]byte
	Type      string `json:"type,omitempty"`
	StringKey string `json:"stringKey,omitempty"`
	HexKey    string `json:"hexKey,omitempty"`
	X         int32  `json:"x,omitempty"`
	Z         int32  `json:"z,omitempty"`
	Y         int32  `json:"y,omitempty"`
	Dimension int32  `json:"dimension,omitempty"`
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

func byteToInt32(i []byte) (int32, error) {
	var out int32
	var err error
	if len(i) != 4 {
		return 0, errors.New("byteToInt32: input must be 4 bytes")
	}
	buf := bytes.NewReader(i[:])
	err = binary.Read(buf, binary.LittleEndian, &out)
	return out, err
}

// TODO: Handle errors
// TODO: use byteToInt32
func keyToCoords(key []byte) (x, z, y int32) {
	var myy byte
	buf := bytes.NewReader(key[0:4])
	_ = binary.Read(buf, binary.LittleEndian, &x)
	buf = bytes.NewReader(key[4:8])
	_ = binary.Read(buf, binary.LittleEndian, &z)
	myy = key[len(key)-1]
	y = int32(myy)
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
	// if the key is ascii, let's just set that as the type and return
	if outkey.StringKey != "" {
		outkey.Type = outkey.StringKey
		return outkey
	}

	// // "~local_player"
	// if bytes.Equal(key, []byte{0x7e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72}[:]) {
	// 	outkey.Type = "player"
	// }
	// // "player_" TODO: double-check this
	// if bytes.Equal(key[:7], []byte{0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f}[:]) {
	// 	outkey.Type = "player"
	// }

	// chunk-based keys
	// legacy terrain keys
	if (len(key) == 9) || (len(key)) == 10 || (len(key) == 13) || (len(key)) == 14 {
		outkey.X, outkey.Z, _ = keyToCoords(key)
		var dimension int32
		if (len(key) == 13) || (len(key) == 14) {
			dimension, _ = byteToInt32(key[8:12])
		}
		outkey.Dimension = dimension
		switch key[len(key)-1] {
		case 0x2d:
			outkey.Type = "data2d"
		case 0x2e:
			outkey.Type = "legacy_data2d"
		// case 0x2f:
		// 0x2f should be a longer key
		// outkey.Type = "terrain"
		case 0x30:
			outkey.Type = "legacy_terrain"
		case 0x31:
			outkey.Type = "block_entity"
		case 0x32:
			outkey.Type = "entity"
		case 0x33:
			outkey.Type = "pending_ticks"
		case 0x34:
			outkey.Type = "block_extra_data"
		case 0x35:
			outkey.Type = "biome_state"
		case 0x36:
			outkey.Type = "0x36-4-byte-value"
		case 0x76:
			outkey.Type = "version"
		default:
			outkey.Dimension = 0
			outkey.Dimension = 0
			outkey.Z = 0
		}
	}
	if (len(key) == 10) || (len(key) == 14) {
		var dimension int32
		outkey.X, outkey.Z, outkey.Y = keyToCoords(key)
		if (len(key) == 13) || (len(key) == 14) {
			dimension, _ = byteToInt32(key[8:12])
		}
		outkey.Dimension = dimension
		switch key[len(key)-2] {
		case 0x2f:
			outkey.Type = "terrain"
		default:
			outkey.Dimension = 0
			outkey.Dimension = 0
			outkey.Z = 0

		}
	}
	if outkey.Type == "" {
		outkey.Type = "unknown"
	}
	return outkey
}
