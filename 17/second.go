package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Coordinates struct {
	x, y int
}

type State struct {
	location       Coordinates
	direction      rune
	directionSteps uint8
}

type PQItem struct {
	state  State
	weight uint64
}

type PriorityQueue []PQItem

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].weight < pq[j].weight
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(PQItem)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]

	return item
}

var directions = map[rune]Coordinates{
	'N': {x: -1, y: 0},
	'S': {x: 1, y: 0},
	'E': {x: 0, y: 1},
	'W': {x: 0, y: -1},
}

var reverse = map[rune]rune{
	'N': 'S',
	'S': 'N',
	'E': 'W',
	'W': 'E',
}

func determineNextLocations(maze [][]uint8, currentPosition PQItem) []PQItem {
	var nextLocations []PQItem
	width := len(maze[0])
	height := len(maze)

	for dir, move := range directions {
		if dir == reverse[currentPosition.state.direction] {
			continue
		}
		x := currentPosition.state.location.x + move.x
		y := currentPosition.state.location.y + move.y
		if x < 0 || x >= height || y < 0 || y >= width {
			continue
		}

		nextLocation := PQItem{
			state: State{location: Coordinates{
				x: currentPosition.state.location.x + move.x,
				y: currentPosition.state.location.y + move.y,
			},
				direction:      dir,
				directionSteps: 1,
			},
			weight: currentPosition.weight,
		}
		if currentPosition.state.direction == dir {
			nextLocation.state.directionSteps = currentPosition.state.directionSteps + 1
		}
		if nextLocation.state.directionSteps > 10 {
			continue
		}
		if nextLocation.state.direction != currentPosition.state.direction && currentPosition.state.directionSteps < 4 {
			continue
		}
		nextLocation.weight += uint64(maze[x][y])
		nextLocations = append(nextLocations, nextLocation)
	}
	return nextLocations
}

func findShortestPath(maze [][]uint8, start, end Coordinates) uint64 {
	visited := make(map[State]uint64)
	var priorityQueue = PriorityQueue{
		PQItem{
			state: State{
				location:       start,
				direction:      'E',
				directionSteps: 0,
			},
			weight: 0,
		},
	}
	heap.Init(&priorityQueue)
	var currentNode PQItem
	for true {
		currentNode = heap.Pop(&priorityQueue).(PQItem)
		if currentNode.state.location == end && currentNode.state.directionSteps >= 4 {
			break
		}
		nextLocations := determineNextLocations(maze, currentNode)
		for _, nextLocation := range nextLocations {
			if weight, found := visited[nextLocation.state]; !found || weight > nextLocation.weight {
				heap.Push(&priorityQueue, nextLocation)
				visited[nextLocation.state] = nextLocation.weight
			}
		}
	}

	fmt.Println(currentNode.state.directionSteps)
	return currentNode.weight
}

func main() {
	inputFile, _ := os.Open("input")
	scanner := bufio.NewScanner(inputFile)

	var maze [][]uint8
	for scanner.Scan() {
		line := scanner.Text()
		splitNumbers := strings.Split(line, "")
		numbers := make([]uint8, len(splitNumbers))
		for ix := 0; ix < len(splitNumbers); ix++ {
			number, _ := strconv.ParseUint(splitNumbers[ix], 10, 8)
			numbers[ix] = uint8(number)
		}
		maze = append(maze, numbers)
	}
	defer func(start time.Time) {
		fmt.Println("took", time.Now().Sub(start))
	}(time.Now())
	fmt.Println(findShortestPath(
		maze,
		Coordinates{x: 0, y: 0},
		Coordinates{x: len(maze) - 1, y: len(maze[0]) - 1},
	))
}
