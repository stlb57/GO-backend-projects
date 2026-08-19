package main

import (
	"fmt"
	"math/rand"
	"os"
	"snake1/internal/game"
	"snake1/internal/terminal"
	"time"

	"golang.org/x/term"
)

const (
	boardWidth  = 60
	boardHeight = 30
)

func randomFood() game.Point {
	return game.Point{
		X: 1 + rand.Intn(boardWidth-2),
		Y: 1 + rand.Intn(boardHeight-2),
	}
}

func main() {
	score := 0
	snake := game.Snake{
		Direction: "up",
		Body: []game.Point{{
			X: boardWidth / 2,
			Y: boardHeight / 2,
		}},
	}

	food := randomFood()

	// //Test snake body
	// snake.Body = append(snake.Body, game.Point{X: 5, Y: 9})
	// snake.Body = append(snake.Body, game.Point{X: 5, Y: 8})
	// snake.Body = append(snake.Body, game.Point{X: 5, Y: 7})
	// snake.Body = append(snake.Body, game.Point{X: 5, Y: 6})
	// snake.Body = append(snake.Body, game.Point{X: 5, Y: 5})
	// snake.Body = append(snake.Body, game.Point{X: 5, Y: 4})

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Could not enter raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	keyChannel := make(chan byte)

gameLoop:
	for {
		fmt.Println(score)
		//Board drawing
		terminal.DrawBoard(boardWidth, boardHeight, snake.Body, food)
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
			quit := game.Move(key, &snake)
			if quit {
				break gameLoop
			}

		default:
		}

		// Primitive collision detection
		if snake.Body[0].X == 0 || snake.Body[0].X == boardWidth-1 || snake.Body[0].Y == 0 || snake.Body[0].Y == boardHeight-1 {
			break gameLoop
		}

		//food eat
		if snake.Body[0].X == food.X && snake.Body[0].Y == food.Y {
			score++

			add_node := game.Point{
				X: snake.Body[len(snake.Body)-1].X,
				Y: snake.Body[len(snake.Body)-1].Y,
			}
			snake.Body = append(snake.Body, add_node)

			//update food location
			for {
				food = randomFood()

				pos := false
				for k := 0; k < len(snake.Body); k++ {
					if snake.Body[k].X == food.X && snake.Body[k].Y == food.Y {
						pos = true
						break
					}
				}

				if !pos {
					break
				}
			}
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

		// Self collision
		for i := 1; i < len(snake.Body); i++ {
			if snake.Body[0].X == snake.Body[i].X && snake.Body[0].Y == snake.Body[i].Y {
				break gameLoop
			}
		}

		time.Sleep(100 * time.Millisecond)

	}
	fmt.Println("GAME OVER")
	fmt.Printf("SCORE : %d", score)
}
