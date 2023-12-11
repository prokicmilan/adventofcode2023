package main

import (
	"bufio"
	"fmt"
	"os"
)

const EMPTY_FIELD = '.'
const GALAXY_FIELD = '#'
const UNIVERSE_GROWTH_FACTOR uint64 = 1000000

type coordinates struct {
	row    int
	column int
}

func countEmptyLinesCrossed(lowerBound, upperBound int, emptyLineIndexes []int) uint64 {
	var emptyLinesCrossed uint64 = 0
	for _, emptyLineIndex := range emptyLineIndexes {
		if emptyLineIndex > lowerBound && emptyLineIndex < upperBound {
			emptyLinesCrossed++
		}
	}

	return emptyLinesCrossed
}

func calculateDistanceBetweenGalaxies(galaxyA, galaxyB coordinates, emptyRowIndexes []int, emptyColumnIndexes []int) uint64 {

	emptyRowsCrossed := countEmptyLinesCrossed(min(galaxyA.row, galaxyB.row), max(galaxyA.row, galaxyB.row), emptyRowIndexes)
	emptyColumnsCrossed := countEmptyLinesCrossed(min(galaxyA.column, galaxyB.column), max(galaxyB.column, galaxyA.column), emptyColumnIndexes)

	smallDistance := galaxyB.row - galaxyA.row + max(galaxyA.column, galaxyB.column) - min(galaxyA.column, galaxyB.column)
	return uint64(smallDistance) + (emptyRowsCrossed+emptyColumnsCrossed)*(UNIVERSE_GROWTH_FACTOR-1)
}

func solve(input *os.File) uint64 {
	scanner := bufio.NewScanner(input)
	var galaxyPositions []coordinates
	var emptyRows []int
	var emptyColumnCandidates []bool

	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(emptyColumnCandidates) == 0 {
			emptyColumnCandidates = make([]bool, len(line))
			for i := 0; i < len(emptyColumnCandidates); i++ {
				emptyColumnCandidates[i] = true
			}
		}

		emptyRow := true
		for ix, character := range line {
			if character == GALAXY_FIELD {
				galaxyPositions = append(galaxyPositions, coordinates{row: row, column: ix})
				emptyRow = false
				emptyColumnCandidates[ix] = false
			}
		}
		if emptyRow {
			emptyRows = append(emptyRows, row)
		}
		row++
	}
	var emptyColumns []int
	for ix := range emptyColumnCandidates {
		if emptyColumnCandidates[ix] {
			emptyColumns = append(emptyColumns, ix)
		}
	}

	var sum uint64 = 0
	for i := 0; i < len(galaxyPositions); i++ {
		for j := i + 1; j < len(galaxyPositions); j++ {
			sum += calculateDistanceBetweenGalaxies(galaxyPositions[i], galaxyPositions[j], emptyRows, emptyColumns)
		}
	}

	return sum
}

func main() {
	input, _ := os.Open("input")
	defer input.Close()
	fmt.Println(solve(input))
}
