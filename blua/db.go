package blua

import (
	"encoding/hex"

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
		hlt := L.NewTable()
		L.RawSet(db, lua.LString("hex_keys"), hlt)
		/*
			clt := L.NewTable()
			L.RawSet(db, lua.LString("chunk_keys"), clt)
		*/
		for _, k := range keys {
			stringkey, hexkey := convertKey(k)
			// string keys
			if stringkey != "" {
				slt.Append(lua.LString(stringkey))
			} else {
				// handle non-string keys here
				_ = hexkey // temp
			}
			// hex keys
			hlt.Append(lua.LString(hexkey))
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
