package blua

import (
	"fmt"
	"testing"

	"github.com/c-bata/go-prompt"
	lua "github.com/yuin/gopher-lua"
)

// I'm new to go-prompt and don't know how to do without auto-suggest yet, so just using their example
func myCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// I'm thinking of just using tests while developing the package, so will
//   try to put stuff here instead of a new executable or mcpetool
// Go will cache test results; run `go test` with `-count=1` to skip caching the interactive input/output
func TestWhatevs(t *testing.T) {
	fmt.Println("Just testing...")
	L := lua.NewState()
	if err := Blua(L); err != nil {
		t.Error("Blua: ", err.Error())
	}
	l := prompt.Input("> ", myCompleter)
	if err := L.DoString(l); err != nil {
		t.Error("DoString: ", err.Error())
	}
}

// go-prompt example
/*
func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func TestGoPrompt(t *testing.T) {
	fmt.Println("Please select table.")
	tbl := prompt.Input("> ", completer)
	fmt.Println("You selected " + tbl)
}
*/
