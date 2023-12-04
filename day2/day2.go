package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Game struct {
	id                     int
	revealedConstellations []CubeConstellation
}

type CubeConstellation struct {
	redCount   int
	blueCount  int
	greenCount int
}

func (game Game) verifyGame(maxCubeConstellation CubeConstellation) bool {
	for _, revealedConstellation := range game.revealedConstellations {
		result := revealedConstellation.verifyConstellation(maxCubeConstellation)
		if !result {
			return false
		}
	}
	return true
}

func (game Game) getMinimumSubset() CubeConstellation {
	maxGreen, maxBlue, maxRed := 0, 0, 0

	for _, constellation := range game.revealedConstellations {
		if constellation.greenCount > 0 && constellation.greenCount > maxGreen {
			maxGreen = constellation.greenCount
		}
		if constellation.redCount > 0 && constellation.redCount > maxRed {
			maxRed = constellation.redCount
		}
		if constellation.blueCount > 0 && constellation.blueCount > maxBlue {
			maxBlue = constellation.blueCount
		}
	}

	return CubeConstellation{
		redCount:   maxRed,
		blueCount:  maxBlue,
		greenCount: maxGreen,
	}
}

func (cubeConstellation CubeConstellation) getPower() int {
	return cubeConstellation.blueCount * cubeConstellation.redCount * cubeConstellation.greenCount
}

func (cubeConstellation CubeConstellation) verifyConstellation(maxCubeConstellation CubeConstellation) bool {
	return cubeConstellation.blueCount <= maxCubeConstellation.blueCount &&
		cubeConstellation.redCount <= maxCubeConstellation.redCount &&
		cubeConstellation.greenCount <= maxCubeConstellation.greenCount
}

func readFile(file string) []string {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	splitLines := strings.Split(lines, "\n")
	return splitLines
}

func main() {
	lines := readFile("day2/input.txt")
	games, err := parseGames(lines)
	if err != nil {
		println(err)
	}
	sum := addUpValidIds(games, CubeConstellation{
		redCount:   12,
		blueCount:  14,
		greenCount: 13,
	})
	fmt.Println("Sum of valid games", sum)

	sum = addUpMinimumSetPowers(games)
	fmt.Println("Sum of minimum set powers", sum)
}

func addUpMinimumSetPowers(games []Game) int {
	sum := 0

	for _, game := range games {
		minimumSubset := game.getMinimumSubset()
		power := minimumSubset.getPower()
		fmt.Println("Game", game.id, "Power", power, "Minsubset", minimumSubset)
		sum += power
	}

	return sum
}

func addUpValidIds(games []Game, constellation CubeConstellation) int {
	counter := 0

	for _, game := range games {
		isValid := game.verifyGame(constellation)
		if isValid {
			counter += game.id
		}
	}

	return counter
}

func parseGames(lines []string) ([]Game, error) {
	var games []Game
	for _, line := range lines {
		game, err := parseGame(line)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}
func parseGame(line string) (Game, error) {
	var game Game

	gameAndSubsetsStrings := strings.Split(line, ": ")
	gameIdString, _ := strings.CutPrefix(gameAndSubsetsStrings[0], "Game ")
	var err error
	game.id, err = strconv.Atoi(gameIdString)
	if err != nil {
		return game, err
	}

	game.revealedConstellations, err = parseSubsets(gameAndSubsetsStrings[1])
	if err != nil {
		return game, err
	}

	return game, nil
}

func parseSubsets(subsetsString string) ([]CubeConstellation, error) {
	var subsets []CubeConstellation

	subsetStrings := strings.Split(subsetsString, "; ")
	for _, subsetString := range subsetStrings {
		subset, err := parseSubset(subsetString)
		if err != nil {
			return subsets, err
		}
		subsets = append(subsets, subset)
	}

	return subsets, nil
}

func parseSubset(subsetString string) (CubeConstellation, error) {
	var subset CubeConstellation

	cubeCountStrings := strings.Split(subsetString, ", ")
	for _, cubeCountString := range cubeCountStrings {
		cubeCount, found := strings.CutSuffix(cubeCountString, " blue")
		if found {
			atoi, err := strconv.Atoi(cubeCount)
			if err != nil {
				return subset, err
			}
			subset.blueCount = atoi
		}
		cubeCount, found = strings.CutSuffix(cubeCountString, " red")
		if found {
			atoi, err := strconv.Atoi(cubeCount)
			if err != nil {
				return subset, err
			}
			subset.redCount = atoi
		}
		cubeCount, found = strings.CutSuffix(cubeCountString, " green")
		if found {
			atoi, err := strconv.Atoi(cubeCount)
			if err != nil {
				return subset, err
			}
			subset.greenCount = atoi
		}
	}

	return subset, nil
}
