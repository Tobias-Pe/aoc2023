package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Race struct {
	time     int
	distance int
}

func (race Race) getCountPossibleDistancesToWin() int {
	counter := 0
	for i := 0; i < race.time; i++ {
		speed := i
		restTime := race.time - i
		distance := speed * restTime
		if distance > race.distance {
			counter++
		}
	}
	return counter
}

func main() {
	start := time.Now()
	lines := readFile("day6/input.txt")
	races, err := parseRaces(lines)
	if err != nil {
		fmt.Println(err)
		return
	}
	multiplicationOfAllPossibilityCounts := 1
	for _, race := range races {
		possibleDistancesToWin := race.getCountPossibleDistancesToWin()
		multiplicationOfAllPossibilityCounts *= possibleDistancesToWin
	}
	fmt.Println("Part1: Multiplication of the number of ways to beat the records", multiplicationOfAllPossibilityCounts)

	race, err := parseRaceWithFixedKerning(lines)
	if err != nil {
		fmt.Println(err)
		return
	}
	possibleDistancesToWin := race.getCountPossibleDistancesToWin()

	fmt.Println("Part2: Multiplication of the number of ways to beat the records", possibleDistancesToWin)

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

func parseRaces(lines []string) ([]Race, error) {
	var races []Race
	times := strings.Fields(lines[0])
	distances := strings.Fields(lines[1])
	for i := 1; i < len(times); i++ {
		time, err := strconv.Atoi(times[i])
		if err != nil {
			return nil, err
		}
		distance, err := strconv.Atoi(distances[i])
		if err != nil {
			return nil, err
		}

		race := Race{
			time:     time,
			distance: distance,
		}
		races = append(races, race)
	}
	return races, nil
}

func parseRaceWithFixedKerning(lines []string) (Race, error) {
	lines[0], _ = strings.CutPrefix(lines[0], "Time:")
	lines[0] = strings.ReplaceAll(lines[0], " ", "")
	lines[1], _ = strings.CutPrefix(lines[1], "Distance:")
	lines[1] = strings.ReplaceAll(lines[1], " ", "")

	time, err := strconv.Atoi(lines[0])
	if err != nil {
		return Race{}, err
	}

	distance, err := strconv.Atoi(lines[1])
	if err != nil {
		return Race{}, err
	}

	return Race{time: time, distance: distance}, nil
}
