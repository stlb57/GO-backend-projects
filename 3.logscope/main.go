package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type logs struct {
	Timestamp string
	Level     string
	Message   string
}

func main() {
	var logEntries []logs
	file := os.Args[1]
	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening file:", err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 3)
		log := logs{
			Timestamp: parts[0],
			Level:     parts[1],
			Message:   parts[2],
		}
		logEntries = append(logEntries, log)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}
