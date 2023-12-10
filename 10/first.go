package main

import (
	"bufio"
	"fmt"
	"os"
)

type coordinates struct {
	x int
	y int
}

type specification struct {
	nextDirection rune
	move          coordinates
}

func readFile(input *os.File) []string {
	var maze []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		maze = append(maze, scanner.Text())
	}

	return maze
}

func findStartingPosition(maze []string) coordinates {
	for x := range maze {
		for y := range maze[x] {
			if maze[x][y] == 'S' {
				return coordinates{x: x, y: y}
			}
		}
	}

	return coordinates{x: -1, y: -1}
}

func checkPosition(maze []string, x, y int, direction rune, pipeSpec map[rune]map[rune]specification) bool {
	if x < 0 || y < 0 || x >= len(maze) || y >= len(maze[x]) {
		return false
	}
	symbol := maze[x][y]
	_, found := pipeSpec[rune(symbol)][direction]

	return found
}

func goToNextPosition(maze []string, currentPosition coordinates, currentDirection rune) specification {
	pipeSpec := map[rune]map[rune]specification{
		'|': {
			'N': specification{
				nextDirection: 'N', move: coordinates{x: -1, y: 0}},
			'S': specification{
				nextDirection: 'S', move: coordinates{x: 1, y: 0}},
		},
		'-': {
			'W': specification{
				nextDirection: 'W', move: coordinates{x: 0, y: -1}},
			'E': specification{
				nextDirection: 'E', move: coordinates{x: 0, y: 1}},
		},
		'L': {
			'S': specification{
				nextDirection: 'E', move: coordinates{x: 0, y: 1}},
			'W': specification{
				nextDirection: 'N', move: coordinates{x: -1, y: 0}},
		},
		'J': {
			'E': specification{
				nextDirection: 'N', move: coordinates{x: -1, y: 0}},
			'S': specification{
				nextDirection: 'W', move: coordinates{x: 0, y: -1}},
		},
		'7': {
			'E': specification{
				nextDirection: 'S', move: coordinates{x: 1, y: 0}},
			'N': specification{
				nextDirection: 'W', move: coordinates{x: 0, y: -1}},
		},
		'F': {
			'N': specification{
				nextDirection: 'E', move: coordinates{x: 0, y: 1}},
			'W': specification{
				nextDirection: 'S', move: coordinates{x: 1, y: 0}},
		},
	}
	currentPipeSymbol := maze[currentPosition.x][currentPosition.y]
	if currentPipeSymbol == 'S' {
		// this is the first step, we need to find the first valid position we can go to
		// we're going to check all sides - up, right, down, left (N, W, S, E)
		directionsToCheck := []specification{
			{nextDirection: 'N', move: coordinates{x: -1, y: 0}},
			{nextDirection: 'E', move: coordinates{x: 0, y: 1}},
			{nextDirection: 'S', move: coordinates{x: 1, y: 0}},
			{nextDirection: 'W', move: coordinates{x: 0, y: -1}},
		}

		for _, directionToCheck := range directionsToCheck {
			newX := currentPosition.x + directionToCheck.move.x
			newY := currentPosition.y + directionToCheck.move.y
			if checkPosition(maze, currentPosition.x+directionToCheck.move.x, currentPosition.y+directionToCheck.move.y, directionToCheck.nextDirection, pipeSpec) {
				return specification{
					nextDirection: directionToCheck.nextDirection,
					move:          coordinates{x: newX, y: newY},
				}
			}
		}
	}

	moveSpec := pipeSpec[rune(currentPipeSymbol)][currentDirection]
	return specification{
		nextDirection: moveSpec.nextDirection,
		move:          coordinates{x: currentPosition.x + moveSpec.move.x, y: currentPosition.y + moveSpec.move.y},
	}
}

func determineLoopLength(maze []string) int {
	loopLength := 0
	startingPosition := findStartingPosition(maze)
	currentPosition := startingPosition

	moveSpec := goToNextPosition(maze, currentPosition, 'A')
	currentDirection := moveSpec.nextDirection
	currentPosition = moveSpec.move
	loopLength++
	for maze[currentPosition.x][currentPosition.y] != 'S' {
		loopLength++
		// fmt.Println("(", currentDirection, ")", currentPosition.x, ":", currentPosition.y, "->", string(maze[currentPosition.x][currentPosition.y]))
		nextPositionSpecification := goToNextPosition(maze, currentPosition, currentDirection)
		currentPosition = nextPositionSpecification.move
		currentDirection = nextPositionSpecification.nextDirection
	}

	return loopLength
}

func main() {
	inputFile, _ := os.Open("input")
	maze := readFile(inputFile)
	inputFile.Close()

	loopLength := determineLoopLength(maze)
	fmt.Println(loopLength)
	fmt.Println(loopLength / 2)
}
