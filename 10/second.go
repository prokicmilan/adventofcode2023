package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

func initializeMazeWithLoop(maze []string) []string {
	var mazeWithLoop []string

	for _, line := range maze {
		var initial string = strings.Repeat(".", len(line))
		mazeWithLoop = append(mazeWithLoop, initial)
	}

	return mazeWithLoop
}

func determineLoop(maze []string) []string {
	mazeWithLoop := initializeMazeWithLoop(maze)

	startingPosition := findStartingPosition(maze)
	currentLine := mazeWithLoop[startingPosition.x]
	mazeWithLoop[startingPosition.x] = currentLine[:startingPosition.y] + "S" + currentLine[startingPosition.y+1:]
	currentPosition := startingPosition

	moveSpec := goToNextPosition(maze, currentPosition, 'A')
	currentDirection := moveSpec.nextDirection
	currentPosition = moveSpec.move
	for maze[currentPosition.x][currentPosition.y] != 'S' {
		currentLine = mazeWithLoop[currentPosition.x]
		mazeWithLoop[currentPosition.x] = currentLine[:currentPosition.y] + string(maze[currentPosition.x][currentPosition.y]) + currentLine[currentPosition.y+1:]
		nextPositionSpecification := goToNextPosition(maze, currentPosition, currentDirection)
		currentPosition = nextPositionSpecification.move
		currentDirection = nextPositionSpecification.nextDirection
	}

	return mazeWithLoop
}

/*
Determines whether the symbol is in loop by counting the number of vertical segments to the right of the symbol.
If the number of vertical segments is odd, the symbol is in loop
L followed by 7 and F followed by J are considered vertical segments (even if they are separated by a number of horizontal segments),
so we need to keep track if any of these show up and are not broken by . | or non-matching segment
*/
func isInLoop(line string, symbolPosition int) bool {
	numberOfVerticalSegments := 0
	verticalSegmentFollowers := map[rune]rune{
		'L': '7',
		'F': 'J',
	}

	previousVerticalSegment := '0'
	for _, character := range line[symbolPosition+1:] {
		if character == '|' {
			numberOfVerticalSegments++
			continue
		}
		if verticalSegmentFollower, ok := verticalSegmentFollowers[previousVerticalSegment]; ok && verticalSegmentFollower == character {
			numberOfVerticalSegments++
		}
		if _, ok := verticalSegmentFollowers[character]; ok {
			previousVerticalSegment = character
			continue
		}
		if character != '-' {
			previousVerticalSegment = '0'
		}
	}

	return numberOfVerticalSegments%2 != 0
}

func countEnclosedTiles(maze []string) int {
	enclosedCount := 0
	for _, line := range maze {
		lineCount := 0
		for symbolPosition, symbol := range line {
			if symbol == '.' && isInLoop(line, symbolPosition) {
				lineCount++
			}
		}
		// fmt.Println(lineCount)
		enclosedCount += lineCount
	}

	return enclosedCount
}

func main() {
	inputFile, _ := os.Open("input")
	maze := readFile(inputFile)
	inputFile.Close()

	mazeWithLoop := determineLoop(maze)
	// not gonna bother with programatically determining what is S - just hardcoding based on input
	startingPosition := findStartingPosition(mazeWithLoop)
	mazeWithLoop[startingPosition.x] = mazeWithLoop[startingPosition.x][:startingPosition.y] + "L" + mazeWithLoop[startingPosition.x][startingPosition.y+1:]
	fmt.Println(countEnclosedTiles(mazeWithLoop))
}
