package main

import (
	"cmp"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

type UnitWithDistance struct {
	previous *UnitWithDistance
	unit     *Unit
	distance int
	visited  bool
}

type Unit struct {
	row      int
	column   int
	isGalaxy bool
}

func (u UnitWithDistance) getNeighbourUnits(units [][]*UnitWithDistance) []*UnitWithDistance {
	var validUnits []*UnitWithDistance

	possibleUnit := Unit{u.unit.row - 1, u.unit.column, false}
	validUnits = checkIfUnvisitedAndAdd(units, possibleUnit, validUnits)
	possibleUnit = Unit{u.unit.row + 1, u.unit.column, false}
	validUnits = checkIfUnvisitedAndAdd(units, possibleUnit, validUnits)
	possibleUnit = Unit{u.unit.row, u.unit.column - 1, false}
	validUnits = checkIfUnvisitedAndAdd(units, possibleUnit, validUnits)
	possibleUnit = Unit{u.unit.row, u.unit.column + 1, false}
	validUnits = checkIfUnvisitedAndAdd(units, possibleUnit, validUnits)

	return validUnits
}

func checkIfUnvisitedAndAdd(units [][]*UnitWithDistance, possibleUnit Unit, validUnits []*UnitWithDistance) []*UnitWithDistance {
	if possibleUnit.isValidCoordinate(len(units)-1, len(units[0])-1) {
		unit := units[possibleUnit.row][possibleUnit.column]
		if !unit.visited {
			validUnits = append(validUnits, unit)
		}
	}
	return validUnits
}

func (u Unit) isValidCoordinate(maxRow, maxColumn int) bool {
	return u.row >= 0 && u.column >= 0 && u.row <= maxRow && u.column <= maxColumn
}

type Universe struct {
	galaxies []*Unit
	data     [][]*Unit
}

func (universe *Universe) getGalaxyPairCombinations() map[*Unit]map[*Unit]bool {
	var galaxyToOtherGalaxies = map[*Unit]map[*Unit]bool{}

	for i := 0; i < len(universe.galaxies)-1; i++ {
		galaxyToOtherGalaxies[universe.galaxies[i]] = map[*Unit]bool{}
		for j := i + 1; j < len(universe.galaxies); j++ {
			galaxyToOtherGalaxies[universe.galaxies[i]][universe.galaxies[j]] = true
		}
	}

	return galaxyToOtherGalaxies
}

func (universe *Universe) expand() {
	for i := 0; i < len(universe.data); i++ {
		rowHasNoGalaxy := true
		for _, unit := range universe.data[i] {
			if unit.isGalaxy {
				rowHasNoGalaxy = false
				break
			}
		}
		if !rowHasNoGalaxy {
			continue
		}
		var tmpData [][]*Unit
		tmpData = append(tmpData, universe.data[:i+1]...)
		tmpData = append(tmpData, universe.data[i:]...)
		universe.data = tmpData
		i++
	}

	for j := 0; j < len(universe.data[0]); j++ {
		columnHasNoGalaxy := true
		for i := 0; i < len(universe.data); i++ {
			if universe.data[i][j].isGalaxy {
				columnHasNoGalaxy = false
				break
			}
		}
		if !columnHasNoGalaxy {
			continue
		}

		for i := 0; i < len(universe.data); i++ {
			var tmpData []*Unit
			tmpData = append(tmpData, universe.data[i][:j+1]...)
			tmpData = append(tmpData, universe.data[i][j:]...)
			universe.data[i] = tmpData
		}
		j++
	}

	universe.fixIndexes()
}

// Attention this breaks all previous pointers and recreates them
func (universe *Universe) fixIndexes() {
	var newGalaxies []*Unit
	for i, row := range universe.data {
		for j, unit := range row {
			fixedNewUnitPtr := &Unit{
				row:      i,
				column:   j,
				isGalaxy: unit.isGalaxy,
			}
			if fixedNewUnitPtr.isGalaxy {
				newGalaxies = append(newGalaxies, fixedNewUnitPtr)
			}
			universe.data[i][j] = fixedNewUnitPtr
		}
	}
	universe.galaxies = newGalaxies
}

func main() {
	start := time.Now()

	lines := readFile("day11/input.txt")
	universe := parseUniverse(lines)
	universe.expand()

	sum := sumShortestPathsStepsBetweenAllGalaxies(universe)
	fmt.Println("[Part1] Sum of steps to all other galaxies", sum)

	fmt.Println("Finished in", time.Since(start))
}

func sumShortestPathsStepsBetweenAllGalaxies(universe Universe) int {
	galaxyPairs := universe.getGalaxyPairCombinations()

	var producerWaitGroup sync.WaitGroup
	galaxyDistanceSumChannel := make(chan int, len(galaxyPairs))

	resultCounter := 0
	for galaxy, targetGalaxies := range galaxyPairs {
		producerWaitGroup.Add(1)
		go func(galaxy *Unit, targetGalaxies map[*Unit]bool) {
			defer producerWaitGroup.Done()
			sum := 0
			allOtherGalaxiesWithDistance := calcShortestPathToOtherGalaxiesSteps(galaxy, universe)
			for _, otherGalaxie := range allOtherGalaxiesWithDistance {
				_, ok := targetGalaxies[otherGalaxie.unit]
				if ok {
					sum += otherGalaxie.distance
				}
			}
			fmt.Println("Result arrived from routine #", resultCounter)
			resultCounter++
			galaxyDistanceSumChannel <- sum
		}(galaxy, targetGalaxies)
	}

	producerWaitGroup.Wait()
	close(galaxyDistanceSumChannel)

	sum := 0
	for galaxyDistanceSum := range galaxyDistanceSumChannel {
		sum += galaxyDistanceSum
	}

	return sum
}

func calcShortestPathToOtherGalaxiesSteps(galaxy *Unit, universe Universe) []*UnitWithDistance {
	var unvisited []*UnitWithDistance
	var matrix [][]*UnitWithDistance
	for _, row := range universe.data {
		matrix = append(matrix, []*UnitWithDistance{})
		for _, unitPtr := range row {
			unitWithDistance := &UnitWithDistance{
				previous: nil,
				unit:     unitPtr,
				distance: math.MaxInt,
			}
			matrix[unitPtr.row] = append(matrix[unitPtr.row], unitWithDistance)
			unvisited = append(unvisited, unitWithDistance)
		}
	}

	start := matrix[galaxy.row][galaxy.column]
	start.distance = 0

	slices.SortFunc(unvisited, func(a, b *UnitWithDistance) int {
		return cmp.Compare(a.distance, b.distance)
	})

	for len(unvisited) > 0 {
		currentUnit := unvisited[0]
		currentUnit.visited = true
		unvisited = unvisited[1:]

		for _, neighbour := range currentUnit.getNeighbourUnits(matrix) {
			if currentUnit.distance+1 < neighbour.distance {
				neighbour.distance = currentUnit.distance + 1
				neighbour.previous = currentUnit
			}
		}

		slices.SortFunc(unvisited, func(a, b *UnitWithDistance) int {
			return cmp.Compare(a.distance, b.distance)
		})
	}

	var galaxyDistances []*UnitWithDistance
	for _, unit := range universe.galaxies {
		galaxyDistances = append(galaxyDistances, matrix[unit.row][unit.column])
	}
	return galaxyDistances
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

func parseUniverse(lines []string) Universe {
	var universe Universe

	for i, line := range lines {
		universe.data = append(universe.data, []*Unit{})
		for j, char := range []rune(line) {
			unitPtr := &Unit{
				row:      i,
				column:   j,
				isGalaxy: char == '#',
			}
			universe.galaxies = append(universe.galaxies, unitPtr)
			universe.data[i] = append(universe.data[i], unitPtr)
		}
	}

	return universe
}
