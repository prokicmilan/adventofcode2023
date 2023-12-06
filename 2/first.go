package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func isDrawPossible(loadedState map[string]int, draw map[string]int) bool {
	for color, drawnNumber := range draw {
		if loadedState[color] < drawnNumber {
			return false
		}
	}

	return true
}

func parseDraw(draw string) map[string]int {
	parsedDraw := map[string]int{
		"red":   0,
		"green": 0,
		"blue":  0,
	}
	for color := range parsedDraw {
		re := regexp.MustCompile(`(\d+) ` + color)
		match := re.FindStringSubmatch(draw)
		if len(match) != 0 {
			parsedDraw[color], _ = strconv.Atoi(match[1])
		}
	}

	return parsedDraw
}

func solve(file *os.File) int {
	scanner := bufio.NewScanner(file)

	loadedState := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}
	sum := 0
	for scanner.Scan() {
		line := scanner.Text()
		drawStartIndex := strings.Index(line, ":")
		gameId, _ := strconv.Atoi(line[len("game "):drawStartIndex])
		gamePossible := true
		for _, gameDraw := range strings.Split(line, ";") {
			parsedDraw := parseDraw(gameDraw)
			if !isDrawPossible(loadedState, parsedDraw) {
				gamePossible = false
			}
		}
		if gamePossible {
			sum += gameId
		}
	}
	return sum
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Println(solve(file))
}
