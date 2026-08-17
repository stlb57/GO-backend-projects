package game

func Move(key byte, player *Point) {
	switch key {
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
