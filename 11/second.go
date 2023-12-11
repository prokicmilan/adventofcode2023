package main

import (
	"bufio"
	"fmt"
	"os"

	"gonum.org/v1/gonum/stat/combin"
)

const EMPTY_FIELD = '.'
const GALAXY_FIELD = '#'
const UNIVERSE_GROWTH_FACTOR uint64 = 1000000

func readUniverse(input *os.File) ([][]rune, []int, []int) {
	var emptyRows map[int]bool = make(map[int]bool)
	var emptyColumnCandidates []bool
	var universe [][]rune

	scanner := bufio.NewScanner(input)

	rowCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		universeRow := make([]rune, len(line))
		if len(emptyColumnCandidates) == 0 {
			emptyColumnCandidates = make([]bool, len(line))
			for i := 0; i < len(line); i++ {
				emptyColumnCandidates[i] = true
			}
		}
		allEmpty := true
		for ix, character := range line {
			if character != EMPTY_FIELD {
				emptyColumnCandidates[ix] = false
				allEmpty = false
			}
			universeRow[ix] = character
		}
		if allEmpty {
			emptyRows[rowCount] = true
		}
		rowCount++
		universe = append(universe, universeRow)
	}
	emptyRowIndexes := make([]int, 0, len(emptyRows))
	for k, v := range emptyRows {
		if v {
			emptyRowIndexes = append(emptyRowIndexes, k)
		}
	}
	var emptyColumnIndexes []int
	for ix := range emptyColumnCandidates {
		if emptyColumnCandidates[ix] {
			emptyColumnIndexes = append(emptyColumnIndexes, ix)
		}
	}

	return universe, emptyRowIndexes, emptyColumnIndexes
}

type coordinates struct {
	row    int
	column int
}

func findGalaxyPositions(universe [][]rune) []coordinates {
	var result []coordinates
	for universeRowIx, universeRow := range universe {
		for universeColIx, universePoint := range universeRow {
			if universePoint == GALAXY_FIELD {
				result = append(result, coordinates{row: universeRowIx, column: universeColIx})
			}
		}
	}

	return result
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

func solve(universe [][]rune, emptyRows []int, emptyColumns []int) uint64 {
	galaxyPositions := findGalaxyPositions(universe)
	var channels []chan uint64 = make([]chan uint64, combin.Binomial(len(galaxyPositions), 2))
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan uint64)
	}

	channelCount := 0
	for i := 0; i < len(galaxyPositions); i++ {
		for j := i + 1; j < len(galaxyPositions); j++ {
			go func(ii, jj, cc int) {
				channels[cc] <- calculateDistanceBetweenGalaxies(galaxyPositions[ii], galaxyPositions[jj], emptyRows, emptyColumns)
			}(i, j, channelCount)
			channelCount++
		}
	}

	var sum uint64 = 0
	for _, channel := range channels {
		distance := <-channel
		// fmt.Println(distance)
		sum += distance
	}

	return sum
}

func main() {
	input, _ := os.Open("input")

	universe, emptyRows, emptyColumns := readUniverse(input)
	input.Close()
	fmt.Println(solve(universe, emptyRows, emptyColumns))
}

func printUniverse(universe [][]rune) {
	for _, universeRow := range universe {
		for _, universeChar := range universeRow {
			fmt.Print(string(universeChar))
		}
		fmt.Println()
	}
}
