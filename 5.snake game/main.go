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

	//Variables Initialisation
	score := 0
	snake := game.Snake{
		Direction: "up",
		Body: []game.Point{{
			X: boardWidth / 2,
			Y: boardHeight / 2,
		}},
	}

	food := game.RandomFood(boardWidth, boardHeight)
	keyChannel := make(chan byte)

	// Initialize the game state
	oldState, setupSuccess := terminal.SetupInput()
	if !setupSuccess {
		return
	}
	defer terminal.RestoreInput(oldState)

gameLoop:
	for {
		fmt.Println(score)

		//Board drawing
		terminal.DrawBoard(boardWidth, boardHeight, snake.Body, food)
		fmt.Println("WASD = Move | Q = Quit")

		// Check for keyboard input and update the snake's direction
		if terminal.Input(keyChannel, &snake) {
			break gameLoop
		}

		// Check if the snake has collided with the wall
		if game.HitWall(snake, boardWidth, boardHeight) {
			break gameLoop
		}

		// Check if the snake has eaten the food and handle growth
		game.EatFood(boardWidth, boardHeight, &snake, &food, &score)

		// Move the snake according to its current direction
		game.RunSnake(&snake)

		// Check if the snake has collided with itself
		if game.HitSelf(snake) {
			break gameLoop
		}

		// Control the speed of the game loop
		time.Sleep(100 * time.Millisecond)

	}

	// Display the final game result
	fmt.Println("GAME OVER")
	fmt.Printf("SCORE : %d", score)
}
