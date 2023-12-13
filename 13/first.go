package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func findReflectionVertical(pattern []string) (int, int) {
	columnIndex := 0
	var columnA string = ""
	var columnB string = ""
	maxReflectionColumn := 0
	maxReflectionSize := -1
	for ; columnIndex < len(pattern[0])-1; columnIndex++ {
		columnA = ""
		columnB = ""
		for rowIndex := 0; rowIndex < len(pattern); rowIndex++ {
			columnA += string(pattern[rowIndex][columnIndex])
			columnB += string(pattern[rowIndex][columnIndex+1])
		}
		if strings.Compare(columnA, columnB) == 0 {
			reflectionSize := 1
			columnA = ""
			columnB = ""
			for columnOffset := 1; columnIndex+columnOffset+1 < len(pattern[0]) && columnIndex-columnOffset >= 0; columnOffset++ {
				for rowIndex := 0; rowIndex < len(pattern); rowIndex++ {
					columnA += string(pattern[rowIndex][columnIndex-columnOffset])
					columnB += string(pattern[rowIndex][columnIndex+columnOffset+1])
				}
				if strings.Compare(columnA, columnB) == 0 {
					reflectionSize++
				} else {
					break
				}
			}
			if reflectionSize > maxReflectionSize {
				// check if it's reaching left or right edge
				if columnIndex+1-reflectionSize == 0 || columnIndex+1+reflectionSize == len(pattern[0]) {
					maxReflectionSize = reflectionSize
					maxReflectionColumn = columnIndex
				}
			}
		}
	}

	return maxReflectionSize, maxReflectionColumn
}

func findReflectionHorizontal(pattern []string) (int, int) {
	rowIndex := 0
	maxReflectionColumn := 0
	maxReflectionSize := -1
	for ; rowIndex < len(pattern)-1; rowIndex++ {
		rowA := pattern[rowIndex]
		rowB := pattern[rowIndex+1]
		if strings.Compare(rowA, rowB) == 0 {
			reflectionSize := 1
			for rowOffset := 1; rowIndex+rowOffset+1 < len(pattern) && rowIndex-rowOffset >= 0; rowOffset++ {
				rowA := pattern[rowIndex-rowOffset]
				rowB := pattern[rowIndex+rowOffset+1]
				if strings.Compare(rowA, rowB) == 0 {
					reflectionSize++
				} else {
					break
				}
			}
			if reflectionSize > maxReflectionSize {
				maxReflectionSize = reflectionSize
				maxReflectionColumn = rowIndex
			}
		}
	}

	return maxReflectionSize, maxReflectionColumn
}

func findPerfectReflection(pattern []string) uint64 {
	verticalReflectionSize, verticalReflectionColumn := findReflectionVertical(pattern)

	// perfect reflection
	// if verticalReflectionColumn+1-verticalReflectionSize == 0 && verticalReflectionColumn+1+verticalReflectionSize == len(pattern[0]) {
	// 	return uint64(verticalReflectionColumn + 1)
	// }
	// reflection reaching left edge
	if verticalReflectionColumn+1-verticalReflectionSize == 0 {
		return uint64(verticalReflectionColumn + 1)
	}
	// reflection reaching right edge
	if verticalReflectionColumn+1+verticalReflectionSize == len(pattern[0]) {
		return uint64(verticalReflectionColumn + 1)
	}

	_, horizontalReflectionColumn := findReflectionHorizontal(pattern)
	return 100 * uint64(horizontalReflectionColumn+1)
}

func main() {
	inputFile, _ := os.Open("input")

	scanner := bufio.NewScanner(inputFile)

	var sum uint64 = 0
	for scanner.Scan() {
		var pattern []string

		line := scanner.Text()
		for len(line) != 0 {
			pattern = append(pattern, line)
			if scanner.Scan() {
				line = scanner.Text()
				continue
			}
			line = ""
		}
		result := findPerfectReflection(pattern)
		sum += result
		fmt.Println(result)
	}

	fmt.Println(sum)
}
