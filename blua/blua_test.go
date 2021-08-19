package blua

import (
	"fmt"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// I'm thinking of just using tests while developing the package, so will
//   try to put stuff here instead of a new executable or mcpetool
// Go will cache test results; run `go test` with `-count=1` to skip caching the interactive input/output
func TestWhatevs(t *testing.T) {
	var l string
	fmt.Println("Just testing...")
	L := lua.NewState()
	if err := Blua(L); err != nil {
		t.Error("Blua: ", err.Error())
	}
	l = `io.write(db[1], "\n")`
	if err := L.DoString(l); err != nil {
		t.Error("DoString: ", err.Error())
	}
	l = "open_world()"
	if err := L.DoString(l); err != nil {
		t.Error("DoString: ", err.Error())
	}
	l = "db.get_keys()"
	if err := L.DoString(l); err != nil {
		t.Error("DoString: ", err.Error())
	}
	l = `io.write(db.keys[1][1])`
	if err := L.DoString(l); err != nil {
		t.Error("DoString: ", err.Error())
	}
}
