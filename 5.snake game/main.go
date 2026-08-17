package main

import (
	"fmt"
	"os"
	"snake1/internal/game"
	"snake1/internal/terminal"

	"golang.org/x/term"
)

func main() {
	player := game.Point{
		X: 10,
		Y: 5,
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Could not enter raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		terminal.DrawBoard(20, 10, player)
		fmt.Println("WASD = Move | Q = Quit")

		var input = make([]byte, 1)
		os.Stdin.Read(input)
		game.Move(input[0], &player)
	}
}
