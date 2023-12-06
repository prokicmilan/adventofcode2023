package main

import (
	"bufio"
	"fmt"
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
	scanner.Scan()
	distanceLine := strings.Split(scanner.Text(), ":")[1]
	re := regexp.MustCompile(`\d+`)
	var times, distances []uint64
	times = parseLine(timeLine, re)
	distances = parseLine(distanceLine, re)
	var numberOfVariations uint64 = 1

	for i := 0; i < len(times); i++ {
		time := times[i]
		recordDistance := distances[i]

		maxTime := time / 2
		maxDistance := maxTime * (time - maxTime)
		var minTime uint64
		if maxDistance > recordDistance {
			var j uint64
			for j = 1; j < maxTime; j++ {
				if j*(time-j) > recordDistance {
					minTime = j
					break
				}
			}
		}
		currentNumberOfVariations := (maxTime - minTime + 1) * 2
		if time%2 == 0 {
			currentNumberOfVariations--
		}
		numberOfVariations *= currentNumberOfVariations
	}

	return numberOfVariations
}

func main() {
	input, _ := os.Open("input")
	defer input.Close()

	fmt.Println(solve(input))
}
