package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type coordinates struct {
	x, y int
}

func runLight(maze [][]rune, start coordinates, direction rune, visited map[coordinates]map[rune]bool) int {
	if _, found := visited[start][direction]; found {
		return 0
	}
	current := start
	energizedTiles := 0
	for true {
		if current.x < 0 || current.x >= len(maze) || current.y < 0 || current.y >= len(maze[0]) {
			break
		}
		if _, found := visited[current]; !found {
			visited[current] = make(map[rune]bool)
			visited[current][direction] = true
			energizedTiles++
		}
		currentCharacter := maze[current.x][current.y]
		switch direction {
		case 'N':
			switch currentCharacter {
			case '-':
				return energizedTiles + runLight(maze, coordinates{x: current.x, y: current.y - 1}, 'W', visited) +
					runLight(maze, coordinates{x: current.x, y: current.y + 1}, 'E', visited)
			case '\\':
				return energizedTiles + runLight(maze, coordinates{x: current.x, y: current.y - 1}, 'W', visited)

			case '/':
				return energizedTiles + runLight(maze, coordinates{x: current.x, y: current.y + 1}, 'E', visited)
			}
			current.x--
		case 'S':
			switch currentCharacter {
			case '-':
				return energizedTiles + runLight(maze, coordinates{x: current.x, y: current.y - 1}, 'W', visited) +
					runLight(maze, coordinates{x: current.x, y: current.y + 1}, 'E', visited)
			case '\\':
				return energizedTiles + runLight(maze, coordinates{x: current.x, y: current.y + 1}, 'E', visited)

			case '/':
				return energizedTiles + runLight(maze, coordinates{x: current.x, y: current.y - 1}, 'W', visited)
			}
			current.x++
		case 'E':
			switch currentCharacter {
			case '|':
				return energizedTiles + runLight(maze, coordinates{x: current.x - 1, y: current.y}, 'N', visited) +
					runLight(maze, coordinates{x: current.x + 1, y: current.y}, 'S', visited)
			case '\\':
				return energizedTiles + runLight(maze, coordinates{x: current.x + 1, y: current.y}, 'S', visited)

			case '/':
				return energizedTiles + runLight(maze, coordinates{x: current.x - 1, y: current.y}, 'N', visited)

			}
			current.y++
		case 'W':
			switch currentCharacter {
			case '|':
				return energizedTiles + runLight(maze, coordinates{x: current.x - 1, y: current.y}, 'N', visited) +
					runLight(maze, coordinates{x: current.x + 1, y: current.y}, 'S', visited)
			case '\\':
				return energizedTiles + runLight(maze, coordinates{x: current.x - 1, y: current.y}, 'N', visited)

			case '/':
				return energizedTiles + runLight(maze, coordinates{x: current.x + 1, y: current.y}, 'S', visited)
			}
			current.y--
		}
	}
	return energizedTiles
}

func consumeChannels(channels []chan int, maxNumberOfEnergizedTiles int) int {
	for _, channel := range channels {
		numberOfEnergizedTiles := <-channel
		if numberOfEnergizedTiles > maxNumberOfEnergizedTiles {
			maxNumberOfEnergizedTiles = numberOfEnergizedTiles
		}
	}
	return maxNumberOfEnergizedTiles
}

func determineNumberOfEnergizedTiles(maze [][]rune, start coordinates, direction rune) int {
	return runLight(maze, start, direction, make(map[coordinates]map[rune]bool))
}

func solve(maze [][]rune) int {
	var maxNumberOfEnergizedTiles int = 0
	channels := make([]chan int, len(maze)*2)
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan int)
	}

	channelIndex := 0
	for i := 0; i < len(maze); i++ {
		go func(ix, channelIx int, maze [][]rune) {
			channels[channelIx] <- determineNumberOfEnergizedTiles(maze, coordinates{x: ix, y: 0}, 'E')
		}(i, channelIndex, maze)
		go func(ix, channelIx int, maze [][]rune) {
			channels[channelIx+1] <- determineNumberOfEnergizedTiles(maze, coordinates{x: ix, y: len(maze[0]) - 1}, 'W')
		}(i, channelIndex, maze)
		channelIndex += 2
	}
	maxNumberOfEnergizedTiles = consumeChannels(channels, maxNumberOfEnergizedTiles)

	channels = make([]chan int, len(maze[0])*2)
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan int)
	}

	channelIndex = 0
	for i := 0; i < len(maze[0]); i++ {
		go func(ix, channelIx int, maze [][]rune) {
			channels[channelIx] <- determineNumberOfEnergizedTiles(maze, coordinates{x: 0, y: ix}, 'S')
		}(i, channelIndex, maze)
		go func(ix, channelIx int, maze [][]rune) {
			channels[channelIx+1] <- determineNumberOfEnergizedTiles(maze, coordinates{x: len(maze) - 1, y: ix}, 'N')
		}(i, channelIndex, maze)
		channelIndex += 2
	}
	maxNumberOfEnergizedTiles = consumeChannels(channels, maxNumberOfEnergizedTiles)

	return maxNumberOfEnergizedTiles
}

func main() {
	inputFile, _ := os.Open("input")
	scanner := bufio.NewScanner(inputFile)
	var maze [][]rune
	for scanner.Scan() {
		line := scanner.Text()
		var runeLine []rune = make([]rune, len(line))
		for ix, character := range line {
			runeLine[ix] = character
		}
		maze = append(maze, runeLine)
	}
	start := time.Now()
	fmt.Println(solve(maze))
	fmt.Println("second took", time.Now().Sub(start))
}
