package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type coordinates struct {
	x, y int
}

type specification struct {
	direction rune
	move      coordinates
}

type KeyType interface {
	coordinates | specification
}

type ConcurrentSafeMap[T KeyType] struct {
	mutex   sync.Mutex
	safeMap map[T]int
}

func (c *ConcurrentSafeMap[T]) CheckAndSet(key T) {
	c.mutex.Lock()
	if _, found := c.safeMap[key]; !found {
		c.safeMap[key]++
	}
	c.mutex.Unlock()
}

func (c *ConcurrentSafeMap[T]) Check(key T) bool {
	c.mutex.Lock()
	_, found := c.safeMap[key]
	c.mutex.Unlock()

	return found
}

func runLight(maze [][]rune, start coordinates, direction rune, energized *ConcurrentSafeMap[coordinates], waitGroup *sync.WaitGroup, visited *ConcurrentSafeMap[specification]) {
	if found := visited.Check(specification{move: start, direction: direction}); found {
		waitGroup.Done()
		return
	}
	current := start
	for true {
		if current.x < 0 || current.x >= len(maze) || current.y < 0 || current.y >= len(maze[0]) {
			break
		}
		energized.CheckAndSet(current)
		visited.CheckAndSet(specification{move: current, direction: direction})
		currentCharacter := maze[current.x][current.y]
		done := false
		switch direction {
		case 'N':
			if currentCharacter == '-' {
				waitGroup.Add(2)
				go runLight(maze, coordinates{x: current.x, y: current.y - 1}, 'W', energized, waitGroup, visited)
				go runLight(maze, coordinates{x: current.x, y: current.y + 1}, 'E', energized, waitGroup, visited)
				done = true
			}
			if currentCharacter == '\\' {
				waitGroup.Add(1)
				go runLight(maze, coordinates{x: current.x, y: current.y - 1}, 'W', energized, waitGroup, visited)
				done = true
			}
			if currentCharacter == '/' {
				waitGroup.Add(1)
				go runLight(maze, coordinates{x: current.x, y: current.y + 1}, 'E', energized, waitGroup, visited)
				done = true
			}
			current.x--
		case 'S':
			if currentCharacter == '-' {
				waitGroup.Add(2)
				go runLight(maze, coordinates{x: current.x, y: current.y - 1}, 'W', energized, waitGroup, visited)
				go runLight(maze, coordinates{x: current.x, y: current.y + 1}, 'E', energized, waitGroup, visited)
				done = true
			}
			if currentCharacter == '\\' {
				waitGroup.Add(1)
				go runLight(maze, coordinates{x: current.x, y: current.y + 1}, 'E', energized, waitGroup, visited)
				done = true
			}
			if currentCharacter == '/' {
				waitGroup.Add(1)
				go runLight(maze, coordinates{x: current.x, y: current.y - 1}, 'W', energized, waitGroup, visited)
				done = true
			}
			current.x++
		case 'E':
			if currentCharacter == '|' {
				waitGroup.Add(2)
				go runLight(maze, coordinates{x: current.x - 1, y: current.y}, 'N', energized, waitGroup, visited)
				go runLight(maze, coordinates{x: current.x + 1, y: current.y}, 'S', energized, waitGroup, visited)
				done = true
			}
			if currentCharacter == '\\' {
				waitGroup.Add(1)
				go runLight(maze, coordinates{x: current.x + 1, y: current.y}, 'S', energized, waitGroup, visited)
				done = true
			}
			if currentCharacter == '/' {
				waitGroup.Add(1)
				go runLight(maze, coordinates{x: current.x - 1, y: current.y}, 'N', energized, waitGroup, visited)
				done = true
			}
			current.y++
		case 'W':
			if currentCharacter == '|' {
				waitGroup.Add(2)
				go runLight(maze, coordinates{x: current.x - 1, y: current.y}, 'N', energized, waitGroup, visited)
				go runLight(maze, coordinates{x: current.x + 1, y: current.y}, 'S', energized, waitGroup, visited)
				done = true
			}
			if currentCharacter == '\\' {
				waitGroup.Add(1)
				go runLight(maze, coordinates{x: current.x - 1, y: current.y}, 'N', energized, waitGroup, visited)
				done = true
			}
			if currentCharacter == '/' {
				waitGroup.Add(1)
				go runLight(maze, coordinates{x: current.x + 1, y: current.y}, 'S', energized, waitGroup, visited)
				done = true
			}
			current.y--
		}
		if done {
			break
		}
	}
	waitGroup.Done()
}

func consumeChannels(channels []chan uint32, maxNumberOfEnergizedTiles uint32) uint32 {
	for _, channel := range channels {
		numberOfEnergizedTiles := <-channel
		if numberOfEnergizedTiles > maxNumberOfEnergizedTiles {
			maxNumberOfEnergizedTiles = numberOfEnergizedTiles
		}
	}
	return maxNumberOfEnergizedTiles
}

func determineNumberOfEnergizedTiles(maze [][]rune, start coordinates, direction rune) uint32 {
	var energized ConcurrentSafeMap[coordinates]
	var visited ConcurrentSafeMap[specification]
	var waitGroup sync.WaitGroup
	energized.safeMap = make(map[coordinates]int)
	visited.safeMap = make(map[specification]int)

	waitGroup.Add(1)
	go runLight(maze, start, direction, &energized, &waitGroup, &visited)

	waitGroup.Wait()

	var sum uint32 = 0
	for _, value := range energized.safeMap {
		sum += uint32(value)
	}
	return sum
}

func solve(maze [][]rune) uint32 {
	var maxNumberOfEnergizedTiles uint32 = 0
	channels := make([]chan uint32, len(maze))
	for i := 0; i < len(maze); i++ {
		channels[i] = make(chan uint32)
	}
	for i := 0; i < len(maze); i++ {
		go func(ix int, maze [][]rune) {
			channels[ix] <- determineNumberOfEnergizedTiles(maze, coordinates{x: ix, y: 0}, 'E')
		}(i, maze)
	}
	maxNumberOfEnergizedTiles = consumeChannels(channels, maxNumberOfEnergizedTiles)

	for i := 0; i < len(maze); i++ {
		go func(ix int, maze [][]rune) {
			channels[ix] <- determineNumberOfEnergizedTiles(maze, coordinates{x: ix, y: len(maze[0]) - 1}, 'W')
		}(i, maze)
	}
	maxNumberOfEnergizedTiles = consumeChannels(channels, maxNumberOfEnergizedTiles)

	channels = make([]chan uint32, len(maze[0]))
	for i := 0; i < len(maze[0]); i++ {
		channels[i] = make(chan uint32)
	}
	for i := 0; i < len(maze[0]); i++ {
		go func(ix int, maze [][]rune) {
			channels[ix] <- determineNumberOfEnergizedTiles(maze, coordinates{x: 0, y: ix}, 'S')
		}(i, maze)
	}
	maxNumberOfEnergizedTiles = consumeChannels(channels, maxNumberOfEnergizedTiles)

	for i := 0; i < len(maze[0]); i++ {
		go func(ix int, maze [][]rune) {
			channels[ix] <- determineNumberOfEnergizedTiles(maze, coordinates{x: len(maze) - 1, y: ix}, 'N')
		}(i, maze)
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
