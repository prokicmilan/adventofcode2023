package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const MULTIPLICATION_COEFFICIENT = 17
const MODULO_COEFFICIENT = 256

type lensInBox struct {
	index       int
	focalLength uint64
}

type box struct {
	freeIndex int
	lenses    map[string]lensInBox
}

func calculateHash(input string, cache map[string]uint64) uint64 {
	if cacheValue, exists := cache[input]; exists {
		return cacheValue
	}
	var total uint64 = 0
	for _, character := range input {
		total += uint64(character)
		total *= MULTIPLICATION_COEFFICIENT
		total %= MODULO_COEFFICIENT
	}
	cache[input] = total
	return total
}

func main() {
	inputFile, _ := os.Open("input")

	scanner := bufio.NewScanner(inputFile)

	scanner.Scan()
	inputLine := scanner.Text()
	inputFile.Close()
	var sum uint64 = 0
	var boxes = make([]box, 256)
	for ix := 0; ix < len(boxes); ix++ {
		boxes[ix].lenses = make(map[string]lensInBox)
	}
	var hashCache = make(map[string]uint64)
	for _, instruction := range strings.Split(inputLine, ",") {
		var label string
		if strings.ContainsRune(instruction, '=') {
			splitInstruction := strings.Split(instruction, "=")
			label = splitInstruction[0]
			focalLength, _ := strconv.ParseUint(splitInstruction[1], 10, 8)
			boxNumber := calculateHash(label, hashCache)
			indexedBox := boxes[boxNumber]
			lensesInBox := indexedBox.lenses
			if foundLens, exists := lensesInBox[label]; exists {
				foundLens.focalLength = focalLength
				boxes[boxNumber].lenses[label] = foundLens
			} else {
				lensesInBox[label] = lensInBox{
					index:       indexedBox.freeIndex,
					focalLength: focalLength,
				}
				boxes[boxNumber].freeIndex++
			}

		} else {
			splitInstruction := strings.Split(instruction, "-")
			label = splitInstruction[0]
			boxNumber := calculateHash(label, hashCache)
			indexedBox := boxes[boxNumber]
			lensesInBox := indexedBox.lenses
			if foundLens, exists := lensesInBox[label]; exists {
				delete(lensesInBox, label)
				boxes[boxNumber].freeIndex--
				for lensLabel, lensState := range lensesInBox {
					if lensState.index > foundLens.index {
						lensState.index--
						boxes[boxNumber].lenses[lensLabel] = lensState
					}
				}
			}
		}
		// fmt.Println(instruction)
		// for boxIx, b := range boxes {
		// 	if b.freeIndex != 0 {
		// 		fmt.Println(boxIx, b)
		// 	}
		// }
		// fmt.Println("--------")
	}
	for boxIx, boxContents := range boxes {
		if boxContents.freeIndex > 0 {
			for _, lens := range boxContents.lenses {
				sum += uint64(boxIx+1) * uint64(lens.index+1) * lens.focalLength
			}
		}
	}

	fmt.Println(sum)
}
