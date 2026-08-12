package main

import (
	"fmt"
	"os"
	"strconv"
	// "errors"
)

func operate(a int, b int, operation func(int, int) int) {
	ans := operation(a, b)
	fmt.Println(ans)
}

func main() {

	//error case 1 : appropriate argument length
	// if len(os.Args)!=3{

	// }

	add := func(a int, b int) int {
		return a + b
	}
	divide := func(a int, b int) int {
		return a / b
	}
	multiply := func(a int, b int) int {
		return a * b
	}
	subtract := func(a int, b int) int {
		return a - b
	}
	var operation func(int, int) int
	a, _ := strconv.Atoi(os.Args[1])
	b, _ := strconv.Atoi(os.Args[2])
	switch os.Args[3] {
	case "add":
		operation = add
	case "subtract":
		operation = subtract
	case "divide":
		operation = divide
	case "multiply":
		operation = multiply
	}
	operate(a, b, operation)

}
