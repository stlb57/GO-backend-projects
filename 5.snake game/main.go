package main

import (
	"fmt"
	"snake1/internal/game"
	"snake1/internal/terminal"
	"time"
)

const (
	boardWidth  = 60
	boardHeight = 30
)

func main() {
	score := 0
	snake := game.Snake{
		Direction: "up",
		Body: []game.Point{{
			X: boardWidth / 2,
			Y: boardHeight / 2,
		}},
	}

	oldState, setupSuccess := terminal.SetupInput()
	if !setupSuccess {
		return
	}
	defer terminal.RestoreInput(oldState)

	food := game.RandomFood(boardWidth, boardHeight)
	keyChannel := make(chan byte)

gameLoop:
	for {
		fmt.Println(score)
		//Board drawing
		terminal.DrawBoard(boardWidth, boardHeight, snake.Body, food)
		fmt.Println("WASD = Move | Q = Quit")

		if terminal.Input(keyChannel, &snake) {
			break gameLoop
		}

		if game.HitWall(snake, boardWidth, boardHeight) {
			break gameLoop
		}

		game.EatFood(boardWidth, boardHeight, &snake, &food, &score)
		game.RunSnake(&snake)

		if game.HitSelf(snake) {
			break gameLoop
		}

		time.Sleep(100 * time.Millisecond)

	}
	fmt.Println("GAME OVER")
	fmt.Printf("SCORE : %d", score)
}
