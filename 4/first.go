package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
)

func determineNumberOfMatches(winningLine string, drawnLine string) uint32 {
	re := regexp.MustCompile(`\d+`)

	winningNumbers := make(map[string]bool)
	for _, winningNumber := range re.FindAllString(winningLine, -1) {
		winningNumbers[winningNumber] = true
	}
	var numberOfMatches uint32 = 0
	for _, drawnNumber := range re.FindAllString(drawnLine, -1) {
		if winningNumbers[drawnNumber] {
			numberOfMatches++
		}
	}

	return numberOfMatches
}

func solve(input *os.File) uint32 {
	scanner := bufio.NewScanner(input)
	var sum uint32 = 0
	for scanner.Scan() {
		line := scanner.Text()

		numbers := strings.Split(line, "|")
		numberOfMatches := determineNumberOfMatches(strings.Split(numbers[0], ":")[1], numbers[1])
		if numberOfMatches >= 1 {
			sum += uint32(math.Pow(2, float64(numberOfMatches-1)))
		}
	}
	return sum
}

func main() {
	input, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	fmt.Println(solve(input))
}
