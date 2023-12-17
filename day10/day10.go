package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Loop struct {
	symbolMap   map[*GridSymbol]bool
	symbolArray []*GridSymbol
	turns       int
}

type Coordinate struct {
	row    int
	column int
}

func (c Coordinate) modify(rowDiff, columnDiff int) Coordinate {
	return Coordinate{
		row:    c.row + rowDiff,
		column: c.column + columnDiff,
	}
}

func (c Coordinate) getSurrounding() []Coordinate {
	var coordinates []Coordinate
	coordinates = append(coordinates, Coordinate{row: c.row + 1, column: c.column})
	coordinates = append(coordinates, Coordinate{row: c.row - 1, column: c.column})
	coordinates = append(coordinates, Coordinate{row: c.row, column: c.column + 1})
	coordinates = append(coordinates, Coordinate{row: c.row, column: c.column - 1})
	return coordinates
}

func (c Coordinate) isValid(mapLength, mapRowLength int) bool {
	return c.row >= 0 && c.row < mapLength && c.column >= 0 && c.column < mapRowLength
}

type GridSymbol struct {
	symbol   string
	location Coordinate
}

// Turndir -1 left, 0 no turn, 1 right turn
func (g GridSymbol) getCoordinateOnOtherSideAndTurnDirComingFrom(staringCoordinate Coordinate) (Coordinate, int, error) {
	if g.location.column == staringCoordinate.column {
		if g.location.row > staringCoordinate.row {
			// Coming from top
			switch g.symbol {
			case "|":
				return Coordinate{staringCoordinate.row + 2, staringCoordinate.column}, 0, nil
			case "L":
				return Coordinate{staringCoordinate.row + 1, staringCoordinate.column + 1}, -1, nil
			case "J":
				return Coordinate{staringCoordinate.row + 1, staringCoordinate.column - 1}, 1, nil
			default:
				return Coordinate{}, 0, fmt.Errorf("unreachable Pipe tried coming from top")
			}
		}

		// Coming from bottom
		switch g.symbol {
		case "|":
			return Coordinate{staringCoordinate.row - 2, staringCoordinate.column}, 0, nil
		case "7":
			return Coordinate{staringCoordinate.row - 1, staringCoordinate.column - 1}, -1, nil
		case "F":
			return Coordinate{staringCoordinate.row - 1, staringCoordinate.column + 1}, 1, nil
		default:
			return Coordinate{}, 0, fmt.Errorf("unreachable Pipe tried coming from bottom")
		}
	}

	if g.location.column > staringCoordinate.column {
		// Coming from left
		switch g.symbol {
		case "-":
			return Coordinate{staringCoordinate.row, staringCoordinate.column + 2}, 0, nil
		case "J":
			return Coordinate{staringCoordinate.row - 1, staringCoordinate.column + 1}, -1, nil
		case "7":
			return Coordinate{staringCoordinate.row + 1, staringCoordinate.column + 1}, 1, nil
		default:
			return Coordinate{}, 0, fmt.Errorf("unreachable Pipe tried coming from left")
		}
	}

	// Coming from right
	switch g.symbol {
	case "-":
		return Coordinate{staringCoordinate.row, staringCoordinate.column - 2}, 0, nil
	case "L":
		return Coordinate{staringCoordinate.row - 1, staringCoordinate.column - 1}, 1, nil
	case "F":
		return Coordinate{staringCoordinate.row + 1, staringCoordinate.column - 1}, -1, nil
	default:
		return Coordinate{}, 0, fmt.Errorf("unreachable Pipe tried coming from right")
	}
}

type Map struct {
	tileGrid [][]*GridSymbol
	start    *GridSymbol
}

func (m Map) getGridSymbol(coordinate Coordinate) *GridSymbol {
	return m.tileGrid[coordinate.row][coordinate.column]
}

func (m Map) getFurthestPositionsStepsOfLoopAndTheLoopItself() (int, Loop, error) {
	var loop Loop
	loop.symbolMap = map[*GridSymbol]bool{}
	loop.symbolMap[m.start] = true
	var symbolArraySecondHalve []*GridSymbol
	loop.symbolArray = append(loop.symbolArray, m.start)

	possibleConnections := m.start.location.getSurrounding()
	var currentlyAtTiles []*GridSymbol
	var comingFromTiles []*GridSymbol
	for i, connection := range possibleConnections {
		currentSymbolPtr := m.getGridSymbol(connection)
		_, turnDir, err := currentSymbolPtr.getCoordinateOnOtherSideAndTurnDirComingFrom(m.start.location)
		if err != nil {
			continue
		}
		currentlyAtTiles = append(currentlyAtTiles, currentSymbolPtr)
		comingFromTiles = append(comingFromTiles, m.start)
		if i == 0 {
			loop.symbolArray = append(loop.symbolArray, currentSymbolPtr)
		} else {
			symbolArraySecondHalve = append([]*GridSymbol{currentSymbolPtr}, symbolArraySecondHalve...)
		}
		loop.symbolMap[currentSymbolPtr] = true
		loop.turns += turnDir
	}

	currentStep := 1
	for !twoAreTheSame(currentlyAtTiles) {
		for i := range currentlyAtTiles {
			goingOverTile := currentlyAtTiles[i]
			comingFromTile := comingFromTiles[i]
			nextTileCoordinate, turnDir, err := goingOverTile.getCoordinateOnOtherSideAndTurnDirComingFrom(comingFromTile.location)
			nextTile := m.getGridSymbol(nextTileCoordinate)
			if err != nil {
				return -1, Loop{}, err
			}
			comingFromTiles[i] = currentlyAtTiles[i]
			currentlyAtTiles[i] = nextTile
			if i == 0 {
				loop.symbolArray = append(loop.symbolArray, nextTile)
			} else {
				symbolArraySecondHalve = append([]*GridSymbol{nextTile}, symbolArraySecondHalve...)
			}
			loop.symbolMap[nextTile] = true
			loop.turns += turnDir
		}
		currentStep++
	}

	symbolArraySecondHalve = append(symbolArraySecondHalve[:0], symbolArraySecondHalve[1:]...)
	loop.symbolArray = append(loop.symbolArray, symbolArraySecondHalve...)

	return currentStep, loop, nil
}

func (m Map) countTilesEnclosedBy(loop Loop) int {
	var visited = map[*GridSymbol]bool{}
	var alreadyInToBeVisitedCoordinates = map[Coordinate]bool{}
	var toBeVisitedCoordinates []Coordinate
	counter := 0

	loopInside := -1
	if loop.turns > 0 {
		loopInside = 1
	}

	prevSymbol := m.start
	for _, currentSymbol := range loop.symbolArray {
		if prevSymbol == currentSymbol {
			continue
		}

		mod := loopInside
		if prevSymbol.location.row == currentSymbol.location.row {
			// Here is some logic bug where u need to swith the > to a < depending on the input
			if prevSymbol.location.column > currentSymbol.location.column {
				mod *= -1
			}
			m.addCoordinate(visited, &alreadyInToBeVisitedCoordinates, &toBeVisitedCoordinates, prevSymbol.location.modify(mod, 0))
			m.addCoordinate(visited, &alreadyInToBeVisitedCoordinates, &toBeVisitedCoordinates, currentSymbol.location.modify(mod, 0))
		} else if prevSymbol.location.column == currentSymbol.location.column {
			// Here is some logic bug where u need to swith the > to a < depending on the input
			if prevSymbol.location.row < currentSymbol.location.row {
				mod *= -1
			}
			m.addCoordinate(visited, &alreadyInToBeVisitedCoordinates, &toBeVisitedCoordinates, prevSymbol.location.modify(0, mod))
			m.addCoordinate(visited, &alreadyInToBeVisitedCoordinates, &toBeVisitedCoordinates, currentSymbol.location.modify(0, mod))
		}

		prevSymbol = currentSymbol
	}

	for len(toBeVisitedCoordinates) > 0 {
		visiting := m.getGridSymbol(toBeVisitedCoordinates[0])
		toBeVisitedCoordinates = append(toBeVisitedCoordinates[:0], toBeVisitedCoordinates[1:]...)

		_, onLoop := loop.symbolMap[visiting]
		if !onLoop {
			counter++
			m.addSurroundingCoordinates(visited, &alreadyInToBeVisitedCoordinates, &toBeVisitedCoordinates, visiting.location)
		}

		visited[visiting] = true
	}

	fmt.Println("Size of grid:", len(m.tileGrid)*len(m.tileGrid[0]), "; Count of detected Tiles:", counter, ";Loop-length:", len(loop.symbolMap))

	return counter
}

func (m Map) addSurroundingCoordinates(visited map[*GridSymbol]bool, alreadyAdded *map[Coordinate]bool, toBeVisitedCoordinates *[]Coordinate, coordinate Coordinate) {
	surrounding := coordinate.getSurrounding()
	for _, coordinate := range surrounding {
		m.addCoordinate(visited, alreadyAdded, toBeVisitedCoordinates, coordinate)
	}
}

func (m Map) addCoordinate(visited map[*GridSymbol]bool, alreadyAdded *map[Coordinate]bool, toBeVisitedCoordinates *[]Coordinate, coordinate Coordinate) {
	insideMap := coordinate.isValid(len(m.tileGrid), len(m.tileGrid[0]))
	if !insideMap {
		return
	}
	gridSymbol := m.getGridSymbol(coordinate)
	_, alreadyVisited := visited[gridSymbol]
	if alreadyVisited {
		return
	}
	_, isAlreadyToBeVisited := (*alreadyAdded)[gridSymbol.location]
	if isAlreadyToBeVisited {
		return
	}
	(*alreadyAdded)[gridSymbol.location] = true
	*toBeVisitedCoordinates = append(*toBeVisitedCoordinates, coordinate)
}

func twoAreTheSame(coordinates []*GridSymbol) bool {
	var coordinateCounter = map[Coordinate]int{}
	for _, coordinate := range coordinates {
		_, ok := coordinateCounter[coordinate.location]
		if ok {
			return true
		}
		coordinateCounter[coordinate.location] = 1
	}
	return false
}

func main() {
	start := time.Now()

	lines := readFile("day10/input.txt")
	pipeMap := parseMap(lines)
	furthestPositionsStepsOfLoop, loop, err := pipeMap.getFurthestPositionsStepsOfLoopAndTheLoopItself()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("[Part1] How many steps along the loop does it take to get from the starting position to the point farthest from the starting position?:\n", furthestPositionsStepsOfLoop)

	fmt.Println("[Part2] How many tiles are enclosed by the loop? Left of Start:\n", pipeMap.countTilesEnclosedBy(loop))

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

func parseMap(lines []string) Map {
	var pipeMap Map
	pipeMap.tileGrid = [][]*GridSymbol{}
	for i, line := range lines {
		pipeMap.tileGrid = append(pipeMap.tileGrid, []*GridSymbol{})
		for j, r := range []rune(line) {
			symbolString := string(r)
			symbolPtr := &GridSymbol{
				symbol: symbolString,
				location: Coordinate{
					row:    i,
					column: j,
				},
			}
			if symbolString == "S" {
				pipeMap.start = symbolPtr
			}
			pipeMap.tileGrid[i] = append(pipeMap.tileGrid[i], symbolPtr)
		}
	}
	return pipeMap
}
