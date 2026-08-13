package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type logs struct {
	Date      string
	Timestamp string
	Level     string
	Message   string
}

func (log logs) Display() {
	fmt.Printf(
		"Date: %s, Time: %s, Level: %s, Message: %s\n",
		log.Date,
		log.Timestamp,
		log.Level,
		log.Message,
	)
}

type Analyzer interface {
	Analyze([]logs) map[string]int
}

type LevelAnalyzer struct{}

func (a LevelAnalyzer) Analyze(logEntries []logs) map[string]int {
	counts := make(map[string]int)

	for _, log := range logEntries {
		counts[log.Level]++
	}

	return counts
}

type ErrorAnalyzer struct{}

func (a ErrorAnalyzer) Analyze(logEntries []logs) map[string]int {
	errors := make(map[string]int)

	for _, log := range logEntries {
		if log.Level == "ERROR" {
			errors[log.Message]++
		}
	}

	return errors
}

func runAnalysis(analyzer Analyzer, logEntries []logs) {
	result := analyzer.Analyze(logEntries)

	for key, count := range result {
		fmt.Printf("%s: %d\n", key, count)
	}
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

		parts := strings.SplitN(line, " ", 4)

		log := logs{
			Date:      parts[0],
			Timestamp: parts[1],
			Level:     parts[2],
			Message:   parts[3],
		}

		logEntries = append(logEntries, log)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading file:", err)
		return
	}

	fmt.Println("LEVEL ANALYSIS")
	levelAnalyzer := LevelAnalyzer{}
	runAnalysis(levelAnalyzer, logEntries)

	fmt.Println("\nERROR ANALYSIS")
	errorAnalyzer := ErrorAnalyzer{}
	runAnalysis(errorAnalyzer, logEntries)
}
