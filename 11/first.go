package main

import (
	"bufio"
	"fmt"
	"os"

	"gonum.org/v1/gonum/stat/combin"
)

const EMPTY_FIELD = '.'
const GALAXY_FIELD = '#'

func expandUniverse(universe [][]rune, emptyRows map[int]bool, emptyColumns []bool) [][]rune {
	var expandedUniverse [][]rune = make([][]rune, len(universe)+len(emptyRows))

	expandedUniverseIndex := 0
	for ix, universeRow := range universe {
		expandedUniverse[expandedUniverseIndex] = universeRow
		expandedUniverseIndex++
		if emptyRows[ix] {
			expandedUniverse[expandedUniverseIndex] = universeRow
			expandedUniverseIndex++
		}
	}
	alreadyAddedColumns := 0
	for columnIndex, columnEmpty := range emptyColumns {
		if !columnEmpty {
			continue
		}

		for universeIndex, universeRow := range expandedUniverse {
			expandedUniverseRow := make([]rune, len(universeRow)+1)
			expandedUniverseIndex := 0
			for ix := 0; ix < len(universeRow); ix++ {
				expandedUniverseRow[expandedUniverseIndex] = universeRow[ix]
				expandedUniverseIndex++
				if ix-alreadyAddedColumns == columnIndex {
					expandedUniverseRow[expandedUniverseIndex] = EMPTY_FIELD
					expandedUniverseIndex++
				}
			}
			expandedUniverse[universeIndex] = expandedUniverseRow
		}
		alreadyAddedColumns++
	}

	return expandedUniverse
}

func readUniverse(input *os.File) [][]rune {
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

	// return universe
	return expandUniverse(universe, emptyRows, emptyColumnCandidates)
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

func calculateDistanceBetweenGalaxies(galaxyA, galaxyB coordinates) uint64 {
	return uint64(galaxyB.row) - uint64(galaxyA.row) + uint64(max(galaxyA.column, galaxyB.column)) - uint64(min(galaxyA.column, galaxyB.column))
}

func solve(universe [][]rune) uint64 {
	galaxyPositions := findGalaxyPositions(universe)
	var channels []chan uint64 = make([]chan uint64, combin.Binomial(len(galaxyPositions), 2))
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan uint64)
	}

	channelCount := 0
	for i := 0; i < len(galaxyPositions); i++ {
		for j := i + 1; j < len(galaxyPositions); j++ {
			go func(ii, jj, cc int) {
				channels[cc] <- calculateDistanceBetweenGalaxies(galaxyPositions[ii], galaxyPositions[jj])
			}(i, j, channelCount)
			channelCount++
		}
	}

	var sum uint64 = 0
	for _, channel := range channels {
		sum += <-channel
	}

	return sum
}

func main() {
	input, _ := os.Open("input")

	universe := readUniverse(input)
	input.Close()
	fmt.Println(solve(universe))
}

func printUniverse(universe [][]rune) {
	for _, universeRow := range universe {
		for _, universeChar := range universeRow {
			fmt.Print(string(universeChar))
		}
		fmt.Println()
	}
}
