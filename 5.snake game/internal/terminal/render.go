package terminal

import (
	"fmt"
	"snake1/internal/game"
)

func DrawBoard(width int, height int, snake []game.Point, food game.Point) {
	fmt.Print("\033[2J\033[H")

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if (i == 0 || i == height-1) && (j == 0 || j == width-1) {
				fmt.Print("+")
			} else if i == 0 || i == height-1 {
				fmt.Print("-")
			} else if j == 0 || j == width-1 {
				fmt.Print("|")
			} else if j == food.X && i == food.Y {
				fmt.Print(":")
			} else {

				//Render snake body

				// if j == player.X-1 && i == player.Y-1 {
				// 	fmt.Print("@")
				// } else {
				// 	fmt.Print(" ")
				// }
				var pos bool = false
				for k := 0; k < len(snake); k++ {
					if snake[k].X == j && snake[k].Y == i {
						if k == 0 {
							fmt.Print("@")
						} else {
							fmt.Print("#")
						}
						pos = true
					}
				}
				if pos == false {
					fmt.Print(" ")
				}

			}
		}
		fmt.Println()
	}
}
