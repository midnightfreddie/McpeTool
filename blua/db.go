package blua

import (
	"encoding/hex"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func dbModule(L *lua.LState) {
	lt := L.NewTable()
	L.SetGlobal("db", lt)
	L.RawSet(lt, lua.LString("get_keys"), L.NewFunction(dbGetKeys))
	// Unsure if I need to define an empty key set
	// L.RawSet(lt, lua.LString("keys"), L.NewTable())
}

func dbGetKeys(L *lua.LState) int {
	keys, err := myWorld.GetKeys()
	if err != nil {
		panic(err)
	}
	dblt := L.GetGlobal("db")
	if db, ok := dblt.(*lua.LTable); ok {
		klt := L.NewTable()
		L.RawSet(db, lua.LString("raw_keys"), klt)
		slt := L.NewTable()
		L.RawSet(db, lua.LString("string_keys"), slt)
		slt.Append(lua.LString("test-delete-me"))
		clt := L.NewTable()
		L.RawSet(db, lua.LString("chunk_keys"), clt)
		for _, k := range keys {
			// string keys
			// FIXME: This is not identifying any string keys
			//  ah, it's actually just a string, not base64-encooded
			if stringkey, hexkey := convertKey(k); stringkey != "" {
				fmt.Println(stringkey)
				slt.Append(lua.LString(stringkey))
			} else {
				// fmt.Println(err)
				/*
					if k[0] != 0 {
						fmt.Println(string(k[:]))
					}
				*/
				// handle non-string keys here
				_ = hexkey // temp
			}
			// raw keys
			kkt := L.NewTable()
			klt.Append(kkt)
			for _, kk := range k {
				kkt.Append(lua.LNumber(kk))
			}
		}
	}
	return 0
}

// copied from api/api.go ConvertKey
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
