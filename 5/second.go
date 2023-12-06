package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func processParsedMap(parsedMap [][]uint64, seeds []uint64) []uint64 {
	var result []uint64
	sort.Slice(parsedMap[:], func(i, j int) bool {
		return parsedMap[i][1] < parsedMap[j][1]
	})

	for i := 0; i < len(seeds); i += 2 {
		seedStart := seeds[i]
		seedEnd := seeds[i+1]
		seedLength := seedEnd - seedStart + 1

		for seedLength > 0 {
			var mappedRangeStart, mappedRangeEnd uint64
			var minimumRangeStartBiggerThanSeedStart int64 = -1
			foundARange := false
			for _, mapRange := range parsedMap {
				rangeStart := mapRange[1]
				rangeEnd := mapRange[1] + mapRange[2] - 1

				if minimumRangeStartBiggerThanSeedStart == -1 && rangeStart > seedStart {
					minimumRangeStartBiggerThanSeedStart = int64(rangeStart)
				}

				if seedStart >= rangeStart && seedStart <= rangeEnd {
					mappedRangeStart = mapRange[0] - rangeStart + seedStart
					foundARange = true
					if rangeEnd >= seedEnd {
						seedLength = 0
						mappedRangeEnd = mapRange[0] - rangeStart + seedEnd
					} else {
						seedStart += rangeEnd - seedStart + 1
						seedLength = seedEnd - seedStart + 1
						mappedRangeEnd = mapRange[0] - rangeStart + rangeEnd
					}
					result = append(result, mappedRangeStart, mappedRangeEnd)
				}
			}
			if !foundARange {
				if minimumRangeStartBiggerThanSeedStart == -1 {
					minimumRangeStartBiggerThanSeedStart = int64(seedEnd)
				}
				result = append(result, seedStart, uint64(minimumRangeStartBiggerThanSeedStart))
				seedLength = seedEnd - uint64(minimumRangeStartBiggerThanSeedStart)
				seedStart = uint64(minimumRangeStartBiggerThanSeedStart)
			}
		}
	}

	return result
}

func solve(input *os.File) uint64 {
	scanner := bufio.NewScanner(input)

	scanner.Scan()
	line := scanner.Text()
	var seeds []uint64
	for _, number := range strings.Split(strings.Split(line, ": ")[1], " ") {
		parsedNumber, _ := strconv.ParseUint(number, 10, 32)
		seeds = append(seeds, parsedNumber)
	}
	var parsedMap [][]uint64
	var mapOfMaps = make([][][]uint64, 6)
	numberOfParsedMaps := 0
	parsingMap := false

	for i := 0; i < len(seeds); i += 2 {
		seeds[i+1] = seeds[i] + seeds[i+1] - 1
	}

	for scanner.Scan() {
		line = scanner.Text()
		if len(line) == 0 {
			continue
		}
		if !strings.Contains(line, "map") {
			// still reading a map, just parse the line
			parsingMap = true
			var parsedRow []uint64
			for _, number := range strings.Split(line, " ") {
				parsedNumber, _ := strconv.ParseUint(number, 10, 32)
				parsedRow = append(parsedRow, parsedNumber)
			}
			parsedMap = append(parsedMap, parsedRow)
		} else if parsingMap {
			// done with reading map, process it
			for i := 0; i < len(seeds); i += 2 {
				processedSeeds := processParsedMap(parsedMap, seeds[i:i+2])
				fmt.Println(processedSeeds)
			}
			seeds = processParsedMap(parsedMap, seeds)
			// fmt.Println(seeds)
			fmt.Println(slices.Min(seeds))
			fmt.Println("------")
			mapOfMaps[numberOfParsedMaps] = parsedMap
			numberOfParsedMaps++
			parsedMap = make([][]uint64, 0)
			parsingMap = false
		}
	}
	seeds = processParsedMap(parsedMap, seeds)
	fmt.Println(seeds)
	fmt.Println(slices.Min(seeds))
	fmt.Println("------")

	return slices.Min(seeds)
}

func main() {
	input, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	solution := solve(input)

	fmt.Println(solution)
}
