package terminal

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func SetupInput() (*term.State, bool) {

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Could not enter raw mode:", err)
		return nil, false
	}

	return oldState, true
}

func RestoreInput(oldState *term.State) {
	term.Restore(int(os.Stdin.Fd()), oldState)
}
