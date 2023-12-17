package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type History struct {
	values []int
}

func (h History) predictNextValue() int {
	var diffHistory History
	for i := 0; i < len(h.values)-1; i++ {
		diffHistory.values = append(diffHistory.values, h.values[i+1]-h.values[i])
	}
	if diffHistory.isAllZeroes() {
		return h.values[len(h.values)-1]
	}
	nextValue := h.values[len(h.values)-1] + diffHistory.predictNextValue()
	return nextValue
}

func (h History) predictPreviousValue() int {
	var diffHistory History
	for i := 0; i < len(h.values)-1; i++ {
		diffHistory.values = append(diffHistory.values, h.values[i+1]-h.values[i])
	}
	if diffHistory.isAllZeroes() {
		return h.values[0]
	}
	prevValue := h.values[0] - diffHistory.predictPreviousValue()
	return prevValue
}

func (h History) isAllZeroes() bool {
	for _, value := range h.values {
		if value != 0 {
			return false
		}
	}
	return true
}

func main() {
	start := time.Now()

	lines := readFile("day9/input.txt")
	histories, err := parseHistories(lines)
	if err != nil {
		fmt.Println(err)
		return
	}

	sum := 0
	for _, history := range histories {
		value := history.predictNextValue()
		sum += value
	}
	fmt.Println("[Part1] Sum of predictions:", sum)

	sum = 0
	for _, history := range histories {
		value := history.predictPreviousValue()
		sum += value
	}
	fmt.Println("[Part2] Sum of previous:", sum)

	fmt.Println("Finished in", time.Since(start))
}

func readFile(file string) []string {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	split := strings.Split(lines, "\n")
	return split
}

func parseHistories(lines []string) ([]History, error) {
	var histories []History

	for _, line := range lines {
		fields := strings.Fields(line)
		numbers, err := toInts(fields)
		if err != nil {
			return nil, err
		}
		histories = append(histories, History{numbers})
	}

	return histories, nil
}

func toInts(fields []string) ([]int, error) {
	var numbers []int
	for _, field := range fields {
		number, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}
	return numbers, nil
}
