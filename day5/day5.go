package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MapEntry struct {
	destination NumberRange
	source      NumberRange
}

type NumberRange struct {
	start  int64
	length int64
}

func (numRange NumberRange) getEnd() int64 {
	return numRange.start + (numRange.length - 1)
}

type Map struct {
	name    string
	entries []MapEntry
}

func main() {
	start := time.Now()
	blocks := readFile("day5/input.txt")
	run(blocks, false)
	fmt.Println("Finished in", time.Since(start))
}

func run(blocks []string, isPart1 bool) {
	allSeedRanges, maps := parseBlocks(blocks, isPart1)
	var producerWaitGroup sync.WaitGroup
	locationResults := make(chan int64, 10)

	seedNrsPerRoutine := int64(10000000)
	for i, seed := range allSeedRanges {
		for splitIndex := 0; int64(splitIndex) <= seed.length/seedNrsPerRoutine; splitIndex++ {
			numRangeStart := seed.start + int64(splitIndex)*seedNrsPerRoutine
			numRangeLength := seedNrsPerRoutine
			splitRange := NumberRange{
				start:  numRangeStart,
				length: numRangeLength,
			}
			if splitRange.getEnd() > seed.getEnd() {
				splitRange.length = seed.getEnd() - splitRange.start // plus minus 1 ???
			}
			producerWaitGroup.Add(1)
			go routineSeedTasksConsumer(&producerWaitGroup, maps, locationResults, splitRange, strconv.Itoa(i)+";"+strconv.Itoa(splitIndex))
		}
	}

	var consumerWaitGroup sync.WaitGroup
	consumerWaitGroup.Add(1)
	var minLocation int64 = math.MaxInt64
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		resultCounter := 0
		for location := range locationResults {
			fmt.Println("Result arrived from routine #", resultCounter)
			resultCounter++
			minLocation = min(minLocation, location)
		}
	}(&consumerWaitGroup)

	producerWaitGroup.Wait()
	close(locationResults)

	consumerWaitGroup.Wait()
	fmt.Println("Min Location:", minLocation)
}

func routineSeedTasksConsumer(wg *sync.WaitGroup, maps []Map, locationResults chan int64, seedTask NumberRange, routineID string) {
	fmt.Println(routineID, "Started")
	defer wg.Done()
	var minLocation int64 = math.MaxInt64
	for i := seedTask.start; i <= seedTask.getEnd(); i++ {
		locationResult := traverseMaps(i, maps)
		minLocation = min(minLocation, locationResult)
	}
	locationResults <- minLocation
	fmt.Println(routineID, "Finished")
}

func readFile(file string) []string {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	return strings.Split(lines, "\n\n")
}

func parseMap(block string) Map {
	var singleMap Map

	lines := strings.Split(block, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			singleMap.name = fields[0]
			continue
		}

		var entry MapEntry
		atoi, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		entry.destination.start = atoi
		atoi, err = strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		entry.source.start = atoi
		atoi, err = strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		entry.source.length = atoi
		entry.destination.length = atoi

		singleMap.entries = append(singleMap.entries, entry)
	}

	return singleMap
}

func traverseMaps(seed int64, maps []Map) int64 {
	currentlyMappedValue := seed
	for _, singleMap := range maps {
		currentlyMappedValue = singleMap.traverse(currentlyMappedValue)
	}
	return currentlyMappedValue
}

func parseBlocks(blocks []string, isPart1 bool) ([]NumberRange, []Map) {
	var seeds []NumberRange
	var maps []Map

	for i, block := range blocks {
		if i == 0 {
			seeds = parseSeeds(block, isPart1)
		} else {
			singleMap := parseMap(block)
			maps = append(maps, singleMap)
		}
	}

	return seeds, maps
}

func parseSeeds(block string, isPart1 bool) []NumberRange {
	var seeds []NumberRange

	fields := strings.Fields(block)
	if isPart1 {
		for _, field := range fields {
			atoi, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				continue
			}
			seeds = append(seeds, NumberRange{
				start:  atoi,
				length: 1,
			})
		}
	} else {
		for i := 1; i+1 < len(fields); i += 2 {
			start, err := strconv.ParseInt(fields[i], 10, 64)
			if err != nil {
				continue
			}
			length, err := strconv.ParseInt(fields[i+1], 10, 64)
			if err != nil {
				continue
			}
			seeds = append(seeds, NumberRange{
				start:  start,
				length: length,
			})
		}
	}

	return seeds
}

func (singleMap Map) traverse(input int64) int64 {
	for _, entry := range singleMap.entries {
		if input >= entry.source.start && input <= entry.source.getEnd() {
			return entry.destination.start + (input - entry.source.start)
		}
	}
	return input
}
