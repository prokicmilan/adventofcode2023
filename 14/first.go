package main

import (
	"bufio"
	"fmt"
	"os"
)

const ROUND_ROCK = 'O'
const SQUARE_ROCK = '#'

func solve(puzzle []string) uint64 {
	var sum uint64 = 0
	coefficients := make([]uint64, len(puzzle[0]))
	for i := 0; i < len(coefficients); i++ {
		coefficients[i] = uint64(len(puzzle))
	}

	for rowIx, line := range puzzle {
		for ix, character := range line {
			if character == ROUND_ROCK {
				sum += coefficients[ix]
				coefficients[ix]--
			}
			if character == SQUARE_ROCK {
				coefficients[ix] = uint64(len(puzzle) - rowIx - 1)
			}
		}
	}

	return sum
}

func main() {
	inputFile, _ := os.Open("input")
	scanner := bufio.NewScanner(inputFile)
	var puzzle []string

	for scanner.Scan() {
		puzzle = append(puzzle, scanner.Text())
	}
	fmt.Println(solve(puzzle))
}
