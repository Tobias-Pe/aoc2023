package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines := readFile("input.txt")

	var calibrationValues []int

	for _, line := range lines {
		value, err := getCalibrationValue(line)
		if err != nil {
			return
		}
		calibrationValues = append(calibrationValues, value)
	}

	sum := 0

	for i, value := range calibrationValues {
		println(value, lines[i])
		sum += value
	}

	println(sum)
}

func readFile(file string) []string {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)

	return strings.Split(lines, "\n")
}

func getCalibrationValue(line string) (int, error) {
	characters := strings.Split(line, "")
	var firstNumber = ""
	var secondNumber = ""

	var stringWithPossibleSpelledNumbers string
	for _, character := range characters {
		stringWithPossibleSpelledNumbers += character
		_, err := strconv.Atoi(character)
		possibleNumber := findNumber(stringWithPossibleSpelledNumbers)
		if err == nil {
			firstNumber, secondNumber = assignNumbers(firstNumber, secondNumber, character)
		} else if len(possibleNumber) != 0 {
			firstNumber, secondNumber = assignNumbers(firstNumber, secondNumber, possibleNumber)
		}
	}

	var calibrationString string

	if len(secondNumber) == 0 {
		calibrationString = firstNumber + firstNumber
	} else {
		calibrationString = firstNumber + secondNumber
	}

	atoi, err := strconv.Atoi(calibrationString)
	if err != nil {
		return -1, err
	}

	return atoi, nil
}

func findNumber(possiblySpelledNumber string) string {
	if strings.HasSuffix(possiblySpelledNumber, "one") {
		return "1"
	}
	if strings.HasSuffix(possiblySpelledNumber, "two") {
		return "2"
	}
	if strings.HasSuffix(possiblySpelledNumber, "three") {
		return "3"
	}
	if strings.HasSuffix(possiblySpelledNumber, "four") {
		return "4"
	}
	if strings.HasSuffix(possiblySpelledNumber, "five") {
		return "5"
	}
	if strings.HasSuffix(possiblySpelledNumber, "six") {
		return "6"
	}
	if strings.HasSuffix(possiblySpelledNumber, "seven") {
		return "7"
	}
	if strings.HasSuffix(possiblySpelledNumber, "eight") {
		return "8"
	}
	if strings.HasSuffix(possiblySpelledNumber, "nine") {
		return "9"
	}
	return ""
}

func assignNumbers(firstNumber string, secondNumber string, character string) (string, string) {
	if len(firstNumber) == 0 {
		firstNumber = character
	} else {
		secondNumber = character
	}
	return firstNumber, secondNumber
}
