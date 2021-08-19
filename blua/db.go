package blua

import (
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
		clt := L.NewTable()
		L.RawSet(db, lua.LString("chunk_keys"), clt)
		for _, k := range keys {
			// string keys
			fmt.Println("hi")
			/*
				if b64Bytes, err := base64.StdEncoding.DecodeString(string(k[:])); err == nil {
					fmt.Println(b64Bytes[0])
					slt.Append(lua.LString(string(b64Bytes[:])))
					fmt.Println(string(b64Bytes[:]))
				} else {
					// handle non-string keys here
				}
			*/
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
