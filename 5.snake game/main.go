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
	snake.Body = append(snake.Body, game.Point{X: 5, Y: 9})
	snake.Body = append(snake.Body, game.Point{X: 5, Y: 8})
	snake.Body = append(snake.Body, game.Point{X: 5, Y: 7})
	snake.Body = append(snake.Body, game.Point{X: 5, Y: 6})
	snake.Body = append(snake.Body, game.Point{X: 5, Y: 5})
	snake.Body = append(snake.Body, game.Point{X: 5, Y: 4})

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Could not enter raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	keyChannel := make(chan byte)

	for {

		//Board drawing
		terminal.DrawBoard(20, 10, snake.Body)
		fmt.Println("WASD = Move | Q = Quit")


		// Concurrent goroutine to check for keypress
		go func() {
			var input = make([]byte, 1)
			os.Stdin.Read(input)
			keyChannel <- input[0]
		}()


		//Check for key press
		select {
		case key := <-keyChannel:
			game.Move(key, &snake)

		default:
		}

		// Primitive collision detection
		if snake.Body[0].X==0 || snake.Body[0].X==19 || snake.Body[0].Y==0 || snake.Body[0].Y==9{
			return 
		}


		//Keep updating the snake body values
        for i := len(snake.Body) - 1; i > 0; i-- {
    snake.Body[i].X = snake.Body[i-1].X
    snake.Body[i].Y = snake.Body[i-1].Y
}

		// Keep the snake running
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
