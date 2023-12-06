package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func processParsedMap(parsedMap [][]uint64, seeds []uint64) []uint64 {
	var result = make([]uint64, len(seeds))
	for ix, seed := range seeds {
		foundNextValue := false
		for _, row := range parsedMap {
			start := row[1]
			end := row[1] + row[2]
			dstStart := row[0]

			if start <= seed && end >= seed {
				result[ix] = dstStart + seed - start
				foundNextValue = true
				break
			}
		}
		if !foundNextValue {
			result[ix] = seed
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
	parsingMap := false

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
			seeds = processParsedMap(parsedMap, seeds)
			parsedMap = make([][]uint64, 0)
			parsingMap = false
		}
	}
	seeds = processParsedMap(parsedMap, seeds)

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
