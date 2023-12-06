package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parseLine(line string, regex *regexp.Regexp) []uint64 {
	var result []uint64
	for _, number := range regex.FindAllString(line, -1) {
		parsedNumber, _ := strconv.Atoi(number)
		result = append(result, uint64(parsedNumber))
	}

	return result
}

func solve(input *os.File) uint64 {
	scanner := bufio.NewScanner(input)

	scanner.Scan()
	timeLine := strings.Split(scanner.Text(), ":")[1]
	timeLine = strings.ReplaceAll(timeLine, " ", "")
	time, _ := strconv.ParseFloat(timeLine, 10)
	scanner.Scan()
	distanceLine := strings.Split(scanner.Text(), ":")[1]
	distanceLine = strings.ReplaceAll(distanceLine, " ", "")
	distance, _ := strconv.ParseFloat(distanceLine, 10)

	maxTime := time / 2
	minTime_a := -(-time + math.Sqrt(math.Pow(time, 2)-4*distance)) / 2
	minTime_b := -(-time - math.Sqrt(math.Pow(time, 2)-4*distance)) / 2
	minTime := math.Ceil(min(minTime_a, minTime_b))

	numberOfVariations := (uint64(maxTime) - uint64(minTime) + 1) * 2
	if uint64(time)%2 == 0 {
		numberOfVariations--
	}

	return numberOfVariations
}

func main() {
	input, _ := os.Open("input")
	defer input.Close()

	fmt.Println(solve(input))
}
