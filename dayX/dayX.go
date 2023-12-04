package main

import (
	"fmt"
	"os"
	"strings"
)

func readFile(file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	splittedLines := strings.Split(lines, "\n")
	for _, line := range splittedLines {
		fmt.Println(">", line)
	}
}

func main() {
	readFile("dayX/input.txt")
}
