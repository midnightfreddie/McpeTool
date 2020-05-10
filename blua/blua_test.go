package blua

import (
	"fmt"
	"testing"
)

// I'm thinking of just using tests while developing the package, so will
//   try to put stuff here instead of a new executable or mcpetool
func TestWhatevs(t *testing.T) {
	fmt.Println("Just testing...")
	// L := lua.NewState()
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
