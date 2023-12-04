package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Scratchcard struct {
	id             int
	winningNumbers []int
	numbers        map[int]bool
}

func (scratchcard Scratchcard) getPoints() int {
	points := 0
	for _, number := range scratchcard.winningNumbers {
		_, ok := scratchcard.numbers[number]
		if !ok {
			continue
		}
		if points == 0 {
			points = 1
		} else {
			points *= 2
		}
	}
	return points
}

func (scratchcard Scratchcard) getMatches() int {
	matches := 0
	for _, number := range scratchcard.winningNumbers {
		_, ok := scratchcard.numbers[number]
		if !ok {
			continue
		}
		matches++
	}
	return matches
}

func main() {
	content := readFile("day4/input.txt")
	scratchcards, err := getScratchcards(content)
	if err != nil {
		fmt.Println(err)
		return
	}

	sum := 0
	for _, scratchcard := range scratchcards {
		sum += scratchcard.getPoints()
	}
	fmt.Println("Points of all scratchcards:", sum)

	fmt.Println("Part2: Total cards:", calculateWinningScratchCards(scratchcards))

}

func calculateWinningScratchCards(scratchcards []Scratchcard) int {
	var counterMap = make(map[int]int)
	for _, originalScratchcard := range scratchcards {
		_, ok := counterMap[originalScratchcard.id]
		if !ok {
			counterMap[originalScratchcard.id] = 0
		}
		counterMap[originalScratchcard.id]++

		for i := originalScratchcard.id + 1; i <= originalScratchcard.id+originalScratchcard.getMatches() &&
			i <= len(scratchcards); i++ {
			_, ok := counterMap[i]
			if !ok {
				counterMap[i] = 0
			}
			counterMap[i] += counterMap[originalScratchcard.id]
		}
	}

	sum := 0
	for _, count := range counterMap {
		sum += count
	}
	return sum
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

func getScratchcards(lines []string) ([]Scratchcard, error) {
	var scratchCards []Scratchcard

	for _, line := range lines {
		var scratchCard Scratchcard

		cardHeaderAndContent := strings.Split(line, ":")

		cardId, _ := strings.CutPrefix(cardHeaderAndContent[0], "Card")
		atoi, err := strconv.Atoi(strings.Fields(cardId)[0])
		if err != nil {
			return nil, err
		}
		scratchCard.id = atoi

		winningNumbersAndNumbers := strings.Split(cardHeaderAndContent[1], "|")
		winningNumbers, err := toIntArray(strings.Fields(winningNumbersAndNumbers[0]))
		if err != nil {
			return nil, err
		}
		scratchCard.winningNumbers = winningNumbers

		numbers, err := toIntArray(strings.Fields(winningNumbersAndNumbers[1]))
		if err != nil {
			return nil, err
		}
		numberSet := make(map[int]bool)
		for _, number := range numbers {
			numberSet[number] = true
		}
		scratchCard.numbers = numberSet

		scratchCards = append(scratchCards, scratchCard)
	}

	return scratchCards, nil
}

func toIntArray(split []string) ([]int, error) {
	var numbers []int
	for _, number := range split {
		atoi, err := strconv.Atoi(number)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, atoi)
	}
	return numbers, nil
}
