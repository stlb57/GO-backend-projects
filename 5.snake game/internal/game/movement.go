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
