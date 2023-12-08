package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readInstructions(scanner *bufio.Scanner) string {
	scanner.Scan()
	return scanner.Text()
}

func readMap(scanner *bufio.Scanner) map[string][2]string {
	var result = make(map[string][2]string)
	// read blank line
	scanner.Scan()
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, " = ")
		source := splitLine[0]
		destinations := strings.Split(splitLine[1], ", ")
		result[source] = [2]string{destinations[0][1:], destinations[1][:len(destinations[1])-1]}
	}

	return result
}

func solve(instructions string, desertMap map[string][2]string) int {
	numberOfSteps := 0
	stepToIndex := map[rune]int{
		'L': 0,
		'R': 1,
	}

	const finalPositionKey = "ZZZ"
	currentPositionKey := "AAA"
	currentPosition := desertMap[currentPositionKey]

	for currentPositionKey != finalPositionKey {
		for _, step := range instructions {
			currentPositionKey = currentPosition[stepToIndex[step]]
			currentPosition = desertMap[currentPositionKey]
			// fmt.Println("currentPosition =", currentPositionKey)
			// fmt.Println("nextPositions =", currentPosition)
			numberOfSteps++
			if currentPositionKey == finalPositionKey {
				break
			}
		}
	}

	return numberOfSteps
}

func main() {
	input, _ := os.Open("input")
	scanner := bufio.NewScanner(input)

	instructions := readInstructions(scanner)
	desertMap := readMap(scanner)

	fmt.Println(solve(instructions, desertMap))

	input.Close()
}
