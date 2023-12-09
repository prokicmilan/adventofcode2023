package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readFile(input *os.File) [][]int64 {
	var result [][]int64

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		var lineNumbers []int64
		for _, num := range strings.Split(line, " ") {
			parsedNumber, _ := strconv.ParseInt(num, 10, 64)
			lineNumbers = append(lineNumbers, parsedNumber)
		}
		result = append(result, lineNumbers)
	}

	return result
}

func areAllZeroes(arr []int64) bool {
	for _, num := range arr {
		if num != 0 {
			return false
		}
	}
	return true
}

func determineNext(line []int64) int64 {
	if areAllZeroes(line) {
		return 0
	}
	var nextArr []int64 = make([]int64, len(line)-1)
	for ix := 1; ix < len(line); ix++ {
		nextArr[ix-1] = line[ix] - line[ix-1]
	}
	return line[0] - determineNext(nextArr)
}

func determineNextNumber(line []int64, channel chan<- int64) {
	channel <- determineNext(line)
}

func solve(input [][]int64) int64 {
	channels := make([]chan int64, len(input))
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan int64)
	}
	var sum int64 = 0
	for ix, line := range input {
		go determineNextNumber(line, channels[ix])
	}
	for _, channel := range channels {
		sum += <-channel
	}

	return sum
}

func main() {
	inputFile, _ := os.Open("input")
	input := readFile(inputFile)
	inputFile.Close()

	fmt.Println(solve(input))
}
