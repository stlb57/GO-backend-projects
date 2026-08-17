package terminal

import (
	"fmt"
	"snake1/internal/game"
)

func DrawBoard(width int, height int, player game.Point) {
	fmt.Print("\033[2J\033[H")

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if (i == 0 || i == height-1) && (j == 0 || j == width-1) {
				fmt.Print("+")
			} else if i == 0 || i == height-1 {
				fmt.Print("-")
			} else if j == 0 || j == width-1 {
				fmt.Print("|")
			} else {
				if j == player.X-1 && i == player.Y-1 {
					fmt.Print("@")
				} else {
					fmt.Print(" ")
				}
			}
		}
		fmt.Println()
	}
}
