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
	length    int
	color     string
}

type Coordinate struct {
	x, y int
}

func calculateArea(coordinates []Coordinate) int {
	var s1, s2 int
	for i := 0; i < len(coordinates)-1; i++ {
		s1 += coordinates[i].x * coordinates[i+1].y
		s2 += coordinates[i].y * coordinates[i+1].x
	}
	return int(math.Abs(float64(s1-s2)) / 2)
}

var directions = map[string]Coordinate{
	"R": {x: 0, y: 1},
	"L": {x: 0, y: -1},
	"U": {x: -1, y: 0},
	"D": {x: 1, y: 0},
}

func processDigPlan(digPlan []DigStep) ([]Coordinate, int) {
	var turningPoints []Coordinate
	var cirumference = 0
	currentPosition := Coordinate{x: 0, y: 0}
	turningPoints = append(turningPoints, currentPosition)
	for _, digStep := range digPlan {
		cirumference += digStep.length
		direction := directions[digStep.direction]
		currentPosition.x += (direction.x * digStep.length)
		currentPosition.y += (direction.y * digStep.length)
		turningPoints = append(turningPoints, currentPosition)
	}

	return turningPoints, cirumference
}

func main() {
	inputFile, _ := os.Open("input")
	scanner := bufio.NewScanner(inputFile)
	var digPlan []DigStep

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, " ")
		parsedLength, _ := strconv.Atoi(splitLine[1])
		step := DigStep{
			direction: splitLine[0],
			length:    parsedLength,
			color:     splitLine[2],
		}
		digPlan = append(digPlan, step)
	}
	turningPoints, circumference := processDigPlan(digPlan)
	fmt.Println(calculateArea(turningPoints) + (circumference / 2) + 1)

	inputFile.Close()
}
