package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func findDigits(line string) int {
	digits := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
	}
	var first, last int
	for _, character := range line {
		stringCharacter := string(character)
		mappedValue := digits[stringCharacter]
		if first == 0 && mappedValue != 0 {
			first = mappedValue
		}
		if mappedValue != 0 {
			last = mappedValue
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
