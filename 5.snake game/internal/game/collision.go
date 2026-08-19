package game

func HitWall(snake Snake, boardWidth int, boardHeight int) bool {
	// Primitive collision detection
	if snake.Body[0].X == 0 || snake.Body[0].X == boardWidth-1 || snake.Body[0].Y == 0 || snake.Body[0].Y == boardHeight-1 {
		return true
	}
	return false
}

func HitSelf(snake Snake) bool {
	// Self collision
	for i := 1; i < len(snake.Body); i++ {
		if snake.Body[0].X == snake.Body[i].X && snake.Body[0].Y == snake.Body[i].Y {
			return true
		}
	}
	return false
}
