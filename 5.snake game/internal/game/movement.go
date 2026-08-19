package game

func Move(key byte, snake *Snake) bool {
	switch key {
	case 'w':
		// snake.Body[0].Y--
		snake.Direction = "up"
	case 's':
		// snake.Body[0].Y++
		snake.Direction = "down"
	case 'a':
		// snake.Body[0].X--
		snake.Direction = "left"
	case 'd':
		// snake.Body[0].X++
		snake.Direction = "right"
	case 'q':
		return true
	}
	return false
}

//return true -> quit requested
//return false -> continue

func RunSnake(snake *Snake) {
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
}
