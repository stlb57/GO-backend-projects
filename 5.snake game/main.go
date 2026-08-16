package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Point struct {
	X int
	Y int
}

func drawBoard(width int, height int, player Point) {
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

func main() {
	player := Point{
		X: 10,
		Y: 5,
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Could not enter raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		drawBoard(20, 10, player)
		fmt.Println("WASD = Move | Q = Quit")

		var input = make([]byte, 1)
		os.Stdin.Read(input)

		switch input[0] {
		case 'w':
			player.Y--
		case 's':
			player.Y++
		case 'a':
			player.X--
		case 'd':
			player.X++
		case 'q':
			return
		}

	}
}
