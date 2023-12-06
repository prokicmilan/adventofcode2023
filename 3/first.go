package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func checkLine(leftBoundary int, rightBoundary int, line *string) bool {
	var indicesToCheck []int
	for i := leftBoundary - 1; i <= rightBoundary; i++ {
		if i >= 0 && i < len(*line) {
			indicesToCheck = append(indicesToCheck, i)
		}
	}
	for _, index := range indicesToCheck {
		if _, err := strconv.Atoi(string((*line)[index])); err != nil && (*line)[index] != '.' {
			return true
		}
	}
	return false
}

func hasNeighbouringSymbols(
	leftBoundary int,
	rightBoundary int,
	currentLine *string,
	lineAbove *string,
	lineBellow *string,
) bool {
	// check left and right
	if leftBoundary-1 >= 0 && (*currentLine)[leftBoundary-1] != '.' {
		return true
	}
	if rightBoundary < len(*currentLine) && (*currentLine)[rightBoundary] != '.' {
		return true
	}
	if lineAbove != nil {
		if checkLine(leftBoundary, rightBoundary, lineAbove) {
			return true
		}
	}
	if lineBellow != nil {
		if checkLine(leftBoundary, rightBoundary, lineBellow) {
			return true
		}
	}
	return false
}

func solve(file *os.File) int {
	scanner := bufio.NewScanner(file)
	var previousLine, currentLine, nextLine *string
	hasMore := scanner.Scan()
	line := scanner.Text()
	currentLine = &line
	re := regexp.MustCompile(`\d+`)
	sum := 0
	for hasMore {
		hasMore = scanner.Scan()
		line := scanner.Text()
		nextLine = &line
		for _, location := range re.FindAllStringIndex(*currentLine, -1) {
			if hasNeighbouringSymbols(
				location[0],
				location[1],
				currentLine,
				previousLine,
				nextLine,
			) {
				number, _ := strconv.Atoi((*currentLine)[location[0]:location[1]])
				sum += number
			}
		}
		previousLine = currentLine
		currentLine = nextLine
	}
	return sum
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Println(solve(file))
}
