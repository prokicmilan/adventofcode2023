package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coordinate struct {
	x, y int
}

var directions = map[string]Coordinate{
	"N": {x: -1, y: 0},
	"S": {x: 1, y: 0},
	"E": {x: 0, y: 1},
	"W": {x: 0, y: -1},
}

func calculateNumberOfPossibleLocations(maze [][]rune, startingLocation Coordinate, numberOfSteps int) int {
	possibleLocations := make(map[int]map[Coordinate]bool)

	stepCounter := 0
	possibleLocations[stepCounter] = make(map[Coordinate]bool)
	possibleLocations[stepCounter][startingLocation] = true

	for true {
		if currentStepPossibleLocations, found := possibleLocations[stepCounter]; !found || len(currentStepPossibleLocations) == 0 {
			stepCounter++
			continue
		}
		if stepCounter == numberOfSteps {
			return len(possibleLocations[stepCounter])
		}

		currentStepPossibleLocations := possibleLocations[stepCounter]
		var currentLocation Coordinate
		for k := range currentStepPossibleLocations {
			currentLocation = k
			break
		}
		delete(currentStepPossibleLocations, currentLocation)
		possibleLocations[stepCounter] = currentStepPossibleLocations

		for _, direction := range directions {
			x := currentLocation.x + direction.x
			y := currentLocation.y + direction.y
			if x < 0 || x >= len(maze) || y < 0 || y >= len(maze[0]) {
				continue
			}
			if maze[x][y] == '#' {
				continue
			}

			nextLocation := Coordinate{
				x: x,
				y: y,
			}
			if _, found := possibleLocations[stepCounter+1]; !found {
				possibleLocations[stepCounter+1] = make(map[Coordinate]bool)
				possibleLocations[stepCounter+1][nextLocation] = true
			} else {
				possibleLocations[stepCounter+1][nextLocation] = true
			}
		}
	}
	return -1
}

func main() {
	inputFile, _ := os.Open("input")

	scanner := bufio.NewScanner(inputFile)

	var maze [][]rune
	var startingLocation Coordinate
	for scanner.Scan() {
		line := scanner.Text()
		var mazeLine = make([]rune, len(line))
		for ix, character := range line {
			if character == 'S' {
				startingLocation = Coordinate{
					x: len(maze),
					y: ix,
				}
			}
			mazeLine[ix] = character
		}
		maze = append(maze, mazeLine)
	}

	fmt.Println(calculateNumberOfPossibleLocations(maze, startingLocation, 64))
}
