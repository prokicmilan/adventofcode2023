package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func findDigits(line string) int {
	digits := map[string]int{
		"1":     1,
		"2":     2,
		"3":     3,
		"4":     4,
		"5":     5,
		"6":     6,
		"7":     7,
		"8":     8,
		"9":     9,
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
	}
	var first, last int
	firstIndex := 9999999999999999
	lastIndex := -1
	for key, value := range digits {
		leftIndex := strings.Index(line, key)
		if leftIndex != -1 && leftIndex < firstIndex {
			firstIndex = leftIndex
			first = value
		}
		rightIndex := strings.LastIndex(line, key)
		if rightIndex != -1 && rightIndex > lastIndex {
			lastIndex = rightIndex
			last = value
		}
	}

	return 10*first + last
}

func solve(file *os.File) int {
	scanner := bufio.NewScanner(file)
	sum := 0
	for scanner.Scan() {
		number := findDigits(scanner.Text())
		sum += number
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
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
