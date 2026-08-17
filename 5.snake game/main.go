package main

import (
	"fmt"
	"os"
	"snake1/internal/game"
	"snake1/internal/terminal"
	"time"

	"golang.org/x/term"
)

func main() {

	snake := game.Snake{
		Direction: "up",
		Body: []game.Point{game.Point{
			X: 10,
			Y: 5,
		}},
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Could not enter raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	keyChannel := make(chan byte)

	for {
		terminal.DrawBoard(20, 10, snake.Body[0])
		fmt.Println("WASD = Move | Q = Quit")

		go func() {
			var input = make([]byte, 1)
			os.Stdin.Read(input)
			keyChannel <- input[0]
		}()

		select {
		case key := <-keyChannel:
			game.Move(key, &snake)

		default:
		}

		switch snake.Direction {
		case "up":
			snake.Body[0].Y--
		case "down":
			snake.Body[0].Y++
		case "left":
			snake.Body[0].X--
		case "right":
			snake.Body[0].X++
		}
		time.Sleep(200 * time.Millisecond)

	}
}
