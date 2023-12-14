package main

import (
	"bufio"
	"fmt"
	"os"
)

const ROUND_ROCK = 'O'
const SQUARE_ROCK = '#'
const EMPTY_SPOT = '.'
const NUMBER_OF_CYCLES = 1000000000

func convert(puzzle []string) [][]rune {
	var result [][]rune = make([][]rune, len(puzzle))
	for ix := range result {
		result[ix] = make([]rune, len(puzzle[0]))
	}
	for rowIx, line := range puzzle {
		for colIx := range line {
			result[rowIx][colIx] = rune(puzzle[rowIx][colIx])
		}
	}

	return result
}

func tilt(puzzle [][]rune, direction rune) [][]rune {

	switch direction {
	case 'N':
		for colIx := 0; colIx < len(puzzle[0]); colIx++ {
			writeIx := 0
			for rowIx := 0; rowIx < len(puzzle); rowIx++ {
				if puzzle[rowIx][colIx] == ROUND_ROCK {
					puzzle[writeIx][colIx] = ROUND_ROCK
					if rowIx != writeIx {
						puzzle[rowIx][colIx] = EMPTY_SPOT
					}
					writeIx++
				}
				if puzzle[rowIx][colIx] == SQUARE_ROCK {
					writeIx = rowIx + 1
				}
			}
		}
	case 'S':
		for colIx := 0; colIx < len(puzzle[0]); colIx++ {
			writeIx := len(puzzle) - 1
			for rowIx := len(puzzle) - 1; rowIx >= 0; rowIx-- {
				if puzzle[rowIx][colIx] == ROUND_ROCK {
					puzzle[writeIx][colIx] = ROUND_ROCK
					if rowIx != writeIx {
						puzzle[rowIx][colIx] = EMPTY_SPOT
					}
					writeIx--
				}
				if puzzle[rowIx][colIx] == SQUARE_ROCK {
					writeIx = rowIx - 1
				}
			}
		}
	case 'W':
		for rowIx := 0; rowIx < len(puzzle); rowIx++ {
			writeIx := 0
			for colIx := 0; colIx < len(puzzle[0]); colIx++ {
				if puzzle[rowIx][colIx] == ROUND_ROCK {
					puzzle[rowIx][writeIx] = ROUND_ROCK
					if colIx != writeIx {
						puzzle[rowIx][colIx] = EMPTY_SPOT
					}
					writeIx++
				}
				if puzzle[rowIx][colIx] == SQUARE_ROCK {
					writeIx = colIx + 1
				}
			}
		}
	case 'E':
		for rowIx := 0; rowIx < len(puzzle); rowIx++ {
			writeIx := len(puzzle[0]) - 1
			for colIx := len(puzzle[0]) - 1; colIx >= 0; colIx-- {
				if puzzle[rowIx][colIx] == ROUND_ROCK {
					puzzle[rowIx][writeIx] = ROUND_ROCK
					if colIx != writeIx {
						puzzle[rowIx][colIx] = EMPTY_SPOT
					}
					writeIx--
				}
				if puzzle[rowIx][colIx] == SQUARE_ROCK {
					writeIx = colIx - 1
				}
			}
		}
	}

	return puzzle
}

func cycle(puzzle [][]rune) [][]rune {
	puzzle = tilt(puzzle, 'N')
	puzzle = tilt(puzzle, 'W')
	puzzle = tilt(puzzle, 'S')
	puzzle = tilt(puzzle, 'E')
	return puzzle
}

func solve(puzzle [][]rune) uint64 {
	var sum uint64 = 0

	for rowIx, line := range puzzle {
		for _, character := range line {
			if character == ROUND_ROCK {
				sum += uint64(len(puzzle) - rowIx)
			}
		}
	}

	return sum
}

func convertToString(puzzle [][]rune) string {
	result := ""
	for _, line := range puzzle {
		result += string(line)
	}

	return result
}

func main() {
	inputFile, _ := os.Open("input")
	scanner := bufio.NewScanner(inputFile)
	var puzzle []string

	for scanner.Scan() {
		puzzle = append(puzzle, scanner.Text())
	}
	puzzleRunes := convert(puzzle)
	// fmt.Println(solve(puzzleRunes))
	puzzleRunes = cycle(puzzleRunes)
	cache := make(map[string]int)
	stringified := convertToString(puzzleRunes)
	cache[stringified] = 1
	cycleCounter := 1
	for true {
		puzzleRunes = cycle(puzzleRunes)
		cycleCounter++
		stringified = convertToString(puzzleRunes)
		if _, found := cache[stringified]; found {
			break
		}
		cache[stringified] = cycleCounter
	}
	cycleStart := cache[stringified]
	cycleLength := cycleCounter - cycleStart
	cycleEnd := cycleStart + cycleLength - 1
	fmt.Println(cycleStart, cycleEnd, cycleLength)
	var cycleNumber int
	for cycleNumber = cycleStart; cycleNumber <= cycleEnd; cycleNumber++ {
		if (NUMBER_OF_CYCLES-cycleNumber)%cycleLength == 0 {
			break
		}
	}
	for i := 0; i < (cycleNumber + cycleLength - cycleCounter); i++ {
		puzzleRunes = cycle(puzzleRunes)
	}

	fmt.Println(solve(puzzleRunes))
}
