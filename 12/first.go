package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const JOKER_SIGN = '?'
const OPERATIONAL = '.'
const DAMAGED = '#'

func isLineInvalid(line string, broken []int) bool {
	var brokenInLine []int

	brokenCount := 0
	for _, character := range line {
		if character == JOKER_SIGN {
			break
		}
		if character == OPERATIONAL {
			if brokenCount > 0 {
				brokenInLine = append(brokenInLine, brokenCount)
				brokenCount = 0
			}
			continue
		}
		brokenCount++
	}

	if len(brokenInLine) == 0 {
		return false
	}
	if len(brokenInLine) > len(broken) {
		return true
	}
	for ix := range brokenInLine {
		if brokenInLine[ix] != broken[ix] {
			return true
		}
	}
	return false
}

func isLineValid(line string, broken []int) bool {
	var brokenInLine []int

	brokenCount := 0
	for _, character := range line {
		if character == JOKER_SIGN {
			break
		}
		if character == OPERATIONAL {
			if brokenCount > 0 {
				brokenInLine = append(brokenInLine, brokenCount)
				brokenCount = 0
			}
			continue
		}
		brokenCount++
	}
	if brokenCount > 0 {
		brokenInLine = append(brokenInLine, brokenCount)
	}

	if len(brokenInLine) != len(broken) {
		return false
	}
	for ix := range brokenInLine {
		if brokenInLine[ix] != broken[ix] {
			return false
		}
	}
	return true
}

func countPossiblePermutations(line string, startingPosition int, broken []int) uint64 {
	if isLineInvalid(line, broken) {
		return 0
	}
	var numberOfPossiblePermutationsDamaged uint64 = 0
	var numberOfPossiblePermutationsOperational uint64 = 0
	for ix, character := range line[startingPosition:] {
		if character == JOKER_SIGN {
			operationalLine := line[:ix+startingPosition] + string(OPERATIONAL) + line[ix+startingPosition+1:]
			damagedLine := line[:ix+startingPosition] + string(DAMAGED) + line[ix+startingPosition+1:]
			numberOfPossiblePermutationsDamaged = countPossiblePermutations(damagedLine, ix+startingPosition+1, broken)
			numberOfPossiblePermutationsOperational = countPossiblePermutations(operationalLine, ix+startingPosition+1, broken)

			return numberOfPossiblePermutationsDamaged + numberOfPossiblePermutationsOperational
		}
	}
	if !isLineValid(line, broken) {
		return 0
	}
	return 1
}

func solve(input *os.File) uint64 {
	var result uint64 = 0

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, " ")
		numbers := strings.Split(splitLine[1], ",")
		broken := make([]int, len(numbers))
		for ix, number := range numbers {
			parsedNumber, _ := strconv.ParseInt(number, 10, 16)
			broken[ix] = int(parsedNumber)

		}
		cnt := countPossiblePermutations(splitLine[0], 0, broken)
		result += cnt
	}

	return result
}

func main() {
	input, _ := os.Open("input")
	defer input.Close()
	fmt.Println(solve(input))
}
