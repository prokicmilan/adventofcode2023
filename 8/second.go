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

func determineStartingPositions(desertMap map[string][2]string) []string {
	var startingPositions []string

	for position := range desertMap {
		if strings.HasSuffix(position, "A") {
			startingPositions = append(startingPositions, position)
		}
	}

	return startingPositions
}

func traverseCount(instructions string, desertMap map[string][2]string, startingPosition string) int {
	numberOfSteps := 0
	stepToIndex := map[rune]int{
		'L': 0,
		'R': 1,
	}

	currentPositionKey := startingPosition
	currentPosition := desertMap[startingPosition]

	for !strings.HasSuffix(currentPositionKey, "Z") {
		for _, step := range instructions {
			currentPositionKey = currentPosition[stepToIndex[step]]
			currentPosition = desertMap[currentPositionKey]
			// fmt.Println("currentPosition =", currentPositionKey)
			// fmt.Println("nextPositions =", currentPosition)
			numberOfSteps++
			if strings.HasSuffix(currentPositionKey, "Z") {
				break
			}
		}
	}

	return numberOfSteps
}

func gcd(a, b uint64) uint64 {
	first := max(a, b)
	second := min(a, b)
	for second != 0 {
		t := second
		second = first % second
		first = t
	}
	return first
}

func lcm(a uint64, b uint64) uint64 {
	return a * b / gcd(a, b)
}

func lcmArray(numbers []uint64) uint64 {
	result := lcm(numbers[0], numbers[1])

	for i := 2; i < len(numbers); i++ {
		result = lcm(result, numbers[i])
	}

	return result
}

func solve(instructions string, desertMap map[string][2]string) uint64 {
	currentPositions := determineStartingPositions(desertMap)

	var numbers []uint64
	for _, position := range currentPositions {
		numbers = append(numbers, uint64(traverseCount(instructions, desertMap, position)))
		// fmt.Println(position, "->", traverseCount(instructions, desertMap, position))
	}

	return lcmArray(numbers)
}

func main() {
	input, _ := os.Open("input")
	scanner := bufio.NewScanner(input)

	instructions := readInstructions(scanner)
	desertMap := readMap(scanner)

	fmt.Println(solve(instructions, desertMap))

	input.Close()
}
