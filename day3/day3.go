package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Coordinate struct {
	i int
	j int
}

func (c Coordinate) isValid(iMax, jMax int) bool {
	return c.i >= 0 && c.j >= 0 && c.i < iMax && c.j < jMax
}

type Number struct {
	startCoordinate Coordinate
	digits          string
}

func (n *Number) append(digit string, coordinate Coordinate) {
	if len(n.digits) == 0 {
		n.startCoordinate = coordinate
	}
	n.digits += digit
}

func (n *Number) getRealNumber() int {
	atoi, err := strconv.Atoi(n.digits)
	if err != nil {
		fmt.Println(err)
	}
	return atoi
}

func (n *Number) getAllCoordinates() []Coordinate {
	var coordinates []Coordinate
	for j := n.startCoordinate.j; j < n.startCoordinate.j+len(n.digits); j++ {
		coordinates = append(coordinates, Coordinate{
			i: n.startCoordinate.i,
			j: j,
		})
	}
	return coordinates
}

type Symbol struct {
	coordinate Coordinate
	symbol     string
}

func (s Symbol) getAdjacentCoordinates(iMax int, jMax int) []Coordinate {
	var coordinates []Coordinate

	originI := s.coordinate.i
	originJ := s.coordinate.j

	// LEFT
	tbAppended := Coordinate{i: originI - 1, j: originJ - 1}
	if tbAppended.isValid(iMax, jMax) {
		coordinates = append(coordinates, tbAppended)
	}
	tbAppended = Coordinate{i: originI - 1, j: originJ}
	if tbAppended.isValid(iMax, jMax) {
		coordinates = append(coordinates, tbAppended)
	}
	tbAppended = Coordinate{i: originI - 1, j: originJ + 1}
	if tbAppended.isValid(iMax, jMax) {
		coordinates = append(coordinates, tbAppended)
	}
	//MID
	tbAppended = Coordinate{i: originI, j: originJ + 1}
	if tbAppended.isValid(iMax, jMax) {
		coordinates = append(coordinates, tbAppended)
	}
	tbAppended = Coordinate{i: originI, j: originJ - 1}
	if tbAppended.isValid(iMax, jMax) {
		coordinates = append(coordinates, tbAppended)
	}
	// RIGHT
	tbAppended = Coordinate{i: originI + 1, j: originJ - 1}
	if tbAppended.isValid(iMax, jMax) {
		coordinates = append(coordinates, tbAppended)
	}
	tbAppended = Coordinate{i: originI + 1, j: originJ}
	if tbAppended.isValid(iMax, jMax) {
		coordinates = append(coordinates, tbAppended)
	}
	tbAppended = Coordinate{i: originI + 1, j: originJ + 1}
	if tbAppended.isValid(iMax, jMax) {
		coordinates = append(coordinates, tbAppended)
	}

	return coordinates
}

type EngineSchematic struct {
	schematic [][]string
	symbols   map[Coordinate]Symbol
	numbers   map[Coordinate]Number
}

func (engineSchematic *EngineSchematic) saveAndResetNumber(number *Number) {
	if len(number.digits) == 0 {
		return
	}

	for _, coordinate := range number.getAllCoordinates() {
		engineSchematic.numbers[coordinate] = *number
	}

	*number = Number{}
}

func (engineSchematic *EngineSchematic) saveSymbol(character string, coordinate Coordinate) {
	engineSchematic.symbols[coordinate] = Symbol{
		coordinate: coordinate,
		symbol:     character,
	}
}

func (engineSchematic *EngineSchematic) populateSymbolsAndNumbers() {
	for i, line := range engineSchematic.schematic {
		var number Number
		for j, character := range line {
			coordinate := Coordinate{i: i, j: j}
			_, err := strconv.Atoi(character)
			if err == nil {
				number.append(character, coordinate)
				continue
			}
			engineSchematic.saveAndResetNumber(&number)
			if character == "." {
				continue
			}
			engineSchematic.saveSymbol(character, coordinate)
		}
		engineSchematic.saveAndResetNumber(&number)
	}
}

func (engineSchematic *EngineSchematic) getPartNumbers() []Number {
	var partNumbersSet = make(map[Number]bool)

	for _, symbol := range engineSchematic.symbols {
		adjacentCoordinates := symbol.getAdjacentCoordinates(len(engineSchematic.schematic), len(engineSchematic.schematic[0]))
		for _, toCheckCoordinate := range adjacentCoordinates {
			number, ok := engineSchematic.numbers[toCheckCoordinate]
			if ok {
				//fmt.Println("Found", number, "next to", symbol)
				partNumbersSet[number] = true
			}
		}
	}

	var partNumbers []Number
	for number, b := range partNumbersSet {
		if b {
			partNumbers = append(partNumbers, number)
		}
	}

	return partNumbers
}

func (engineSchematic *EngineSchematic) getGearRatios() []int {
	var gearRatios []int

	for _, symbol := range engineSchematic.symbols {
		// only gears
		if symbol.symbol != "*" {
			continue
		}
		var partNumbersSet = make(map[Number]bool)
		adjacentCoordinates := symbol.getAdjacentCoordinates(len(engineSchematic.schematic), len(engineSchematic.schematic[0]))
		for _, toCheckCoordinate := range adjacentCoordinates {
			number, ok := engineSchematic.numbers[toCheckCoordinate]
			if ok {
				//fmt.Println("Found", number, "next to", symbol)
				partNumbersSet[number] = true
			}
		}
		if len(partNumbersSet) != 2 {
			continue
		}
		gearRatio := 1
		for number := range partNumbersSet {
			gearRatio *= number.getRealNumber()
		}

		gearRatios = append(gearRatios, gearRatio)
	}

	return gearRatios
}

func main() {
	start := time.Now()
	lines := readFile("day3/input.txt")
	engineSchematic := generateEngineSchematic(lines)
	engineSchematic.populateSymbolsAndNumbers()
	numbers := engineSchematic.getPartNumbers()

	fmt.Println("Sum of part numbers", sumPartNumbers(numbers))

	ratios := engineSchematic.getGearRatios()
	fmt.Println("Sum of part numbers", sumGearRatios(ratios))
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
	return strings.Split(lines, "\n")
}

func generateEngineSchematic(lines []string) EngineSchematic {
	var schematic [][]string
	for _, line := range lines {
		characters := strings.Split(line, "")
		schematic = append(schematic, characters)
	}

	return EngineSchematic{
		schematic: schematic,
		symbols:   make(map[Coordinate]Symbol),
		numbers:   make(map[Coordinate]Number),
	}
}

func sumPartNumbers(numbers []Number) int {
	sum := 0
	for _, number := range numbers {
		sum += number.getRealNumber()
	}
	return sum
}

func sumGearRatios(ratios []int) int {
	sum := 0
	for _, ratio := range ratios {
		sum += ratio
	}
	return sum
}
