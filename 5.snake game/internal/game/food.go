package game

import (
	"math/rand"
)

func RandomFood(boardWidth int, boardHeight int) Point {
	return Point{
		X: 1 + rand.Intn(boardWidth-2),
		Y: 1 + rand.Intn(boardHeight-2),
	}
}

func EatFood(boardWidth int, boardHeight int, snake *Snake, food *Point, score *int) {
	//food eat
	if snake.Body[0].X == food.X && snake.Body[0].Y == food.Y {
		(*score)++

		add_node := Point{
			X: snake.Body[len(snake.Body)-1].X,
			Y: snake.Body[len(snake.Body)-1].Y,
		}
		snake.Body = append(snake.Body, add_node)

		//update food location
		for {
			*food = RandomFood(boardWidth, boardHeight)

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
}
