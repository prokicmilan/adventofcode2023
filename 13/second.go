package main

import (
	"bufio"
	"fmt"
	"os"
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
		fixedSmudge := false
		lineDiff := 0
		for colIx := 0; colIx < len(columnA); colIx++ {
			if columnA[colIx] != columnB[colIx] {
				lineDiff++
			}
		}
		if lineDiff == 1 {
			lineDiff = 0
			fixedSmudge = true
		}

		if lineDiff == 0 {
			reflectionSize := 1
			for columnOffset := 1; columnIndex+columnOffset+1 < len(pattern[0]) && columnIndex-columnOffset >= 0; columnOffset++ {
				columnA = ""
				columnB = ""
				for rowIndex := 0; rowIndex < len(pattern); rowIndex++ {
					columnA += string(pattern[rowIndex][columnIndex-columnOffset])
					columnB += string(pattern[rowIndex][columnIndex+columnOffset+1])
				}
				lineDiff = 0
				for ix := 0; ix < len(columnA); ix++ {
					if columnA[ix] != columnB[ix] {
						lineDiff++
					}
				}
				if lineDiff == 0 {
					reflectionSize++
				} else if lineDiff == 1 && !fixedSmudge {
					reflectionSize++
					fixedSmudge = true
				} else {
					break
				}
			}
			if reflectionSize > maxReflectionSize && fixedSmudge {
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
		fixedSmudge := false
		lineDiff := 0

		for rowIx := 0; rowIx < len(rowA); rowIx++ {
			if rowA[rowIx] != rowB[rowIx] {
				lineDiff++
			}
		}
		if lineDiff == 1 {
			lineDiff = 0
			fixedSmudge = true
		}

		if lineDiff == 0 {
			reflectionSize := 1
			for rowOffset := 1; rowIndex+rowOffset+1 < len(pattern) && rowIndex-rowOffset >= 0; rowOffset++ {
				rowA := pattern[rowIndex-rowOffset]
				rowB := pattern[rowIndex+rowOffset+1]
				lineDiff = 0
				for ix := 0; ix < len(rowA); ix++ {
					if rowA[ix] != rowB[ix] {
						lineDiff++
					}
				}
				if lineDiff == 0 {
					reflectionSize++
				} else if lineDiff == 1 && !fixedSmudge {
					reflectionSize++
					fixedSmudge = true
				} else {
					break
				}
			}
			if reflectionSize > maxReflectionSize && fixedSmudge {
				// check if it's reaching top or bottom edge
				if rowIndex+1-reflectionSize == 0 || rowIndex+1+reflectionSize == len(pattern) {
					maxReflectionSize = reflectionSize
					maxReflectionColumn = rowIndex
				}
			}
		}
	}

	return maxReflectionSize, maxReflectionColumn
}

func findPerfectReflection(pattern []string) uint64 {
	verticalReflectionSize, verticalReflectionColumn := findReflectionVertical(pattern)
	_, horizontalReflectionColumn := findReflectionHorizontal(pattern)

	// reflection reaching left edge
	if verticalReflectionColumn+1-verticalReflectionSize == 0 {
		return uint64(verticalReflectionColumn + 1)
	}
	// reflection reaching right edge
	if verticalReflectionColumn+1+verticalReflectionSize == len(pattern[0]) {
		return uint64(verticalReflectionColumn + 1)
	}

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
	}

	fmt.Println(sum)
}
