package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const MULTIPLICATION_COEFFICIENT = 17
const MODULO_COEFFICIENT = 256

func calculateHash(input string) uint64 {
	var total uint64 = 0
	for _, character := range input {
		total += uint64(character)
		total *= MULTIPLICATION_COEFFICIENT
		total %= MODULO_COEFFICIENT
	}
	return total
}

func main() {
	inputFile, _ := os.Open("input")

	scanner := bufio.NewScanner(inputFile)

	scanner.Scan()
	inputLine := scanner.Text()
	inputFile.Close()
	var sum uint64 = 0
	for _, instruction := range strings.Split(inputLine, ",") {
		sum += calculateHash(instruction)
	}

	fmt.Println(sum)
}
