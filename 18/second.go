package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type DigStep struct {
	direction string
	length    uint64
}

type Coordinate struct {
	x, y int64
}

func calculateArea(coordinates []Coordinate) uint64 {
	var s1, s2 int64
	for i := 0; i < len(coordinates)-1; i++ {
		s1 += coordinates[i].x * coordinates[i+1].y
		s2 += coordinates[i].y * coordinates[i+1].x
	}
	return uint64(math.Abs(float64(s1-s2)) / 2)
}

var directions = map[string]Coordinate{
	"0": {x: 0, y: 1},
	"1": {x: 1, y: 0},
	"2": {x: 0, y: -1},
	"3": {x: -1, y: 0},
}

func processDigPlan(digPlan []DigStep) ([]Coordinate, uint64) {
	var turningPoints []Coordinate
	var circumference uint64
	currentPosition := Coordinate{x: 0, y: 0}
	turningPoints = append(turningPoints, currentPosition)
	for _, digStep := range digPlan {
		circumference += digStep.length
		direction := directions[digStep.direction]
		currentPosition.x += (direction.x * int64(digStep.length))
		currentPosition.y += (direction.y * int64(digStep.length))
		turningPoints = append(turningPoints, currentPosition)
	}

	return turningPoints, circumference
}

func parseDirectionAndLength(color string) (uint64, string) {
	hex := color[2 : len(color)-2]
	dir := color[len(color)-2 : len(color)-1]
	parsedNumber, _ := strconv.ParseUint(hex, 16, 64)
	return parsedNumber, dir
}

func main() {
	inputFile, _ := os.Open("input")
	scanner := bufio.NewScanner(inputFile)
	var digPlan []DigStep

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, " ")
		length, direction := parseDirectionAndLength(splitLine[2])
		step := DigStep{
			length:    length,
			direction: direction,
		}
		digPlan = append(digPlan, step)
	}
	fmt.Println(digPlan)
	turningPoints, circumference := processDigPlan(digPlan)
	fmt.Println(calculateArea(turningPoints) + (circumference / 2) + 1)

	inputFile.Close()
}
