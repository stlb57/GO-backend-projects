package terminal

import (
	"os"
	"snake1/internal/game"
)

func Input(keyChannel chan byte, snake *game.Snake) bool {
	// Concurrent goroutine to check for keypress
	go func() {
		var input = make([]byte, 1)
		os.Stdin.Read(input)
		keyChannel <- input[0]
	}()

	//Check for key press
	select {
	case key := <-keyChannel:
		quit := game.Move(key, snake)
		if quit {
			return true
		}

	default:
	}
	return false
}
