package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func determineNumberOfMatches(winningLine string, drawnLine string) int {
	re := regexp.MustCompile(`\d+`)

	winningNumbers := make(map[string]bool)
	for _, winningNumber := range re.FindAllString(winningLine, -1) {
		winningNumbers[winningNumber] = true
	}
	var numberOfMatches int = 0
	for _, drawnNumber := range re.FindAllString(drawnLine, -1) {
		if winningNumbers[drawnNumber] {
			numberOfMatches++
		}
	}

	return numberOfMatches
}

func solve(input *os.File) int {
	scanner := bufio.NewScanner(input)
	gameId := 1
	tickets := make(map[int]int)
	for scanner.Scan() {
		line := scanner.Text()
		tickets[gameId]++
		numbers := strings.Split(line, "|")
		numberOfMatches := determineNumberOfMatches(strings.Split(numbers[0], ":")[1], numbers[1])
		for i := 1; i <= numberOfMatches; i++ {
			tickets[gameId+i] += tickets[gameId]
		}

		gameId++
	}
	sum := 0
	for _, numberOfTickets := range tickets {
		sum += numberOfTickets
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
