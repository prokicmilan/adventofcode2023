package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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

func determineMinimumLoadout(currentMinimumLoadout map[string]int, draw map[string]int) map[string]int {
	minimumLoadoutNecessary := make(map[string]int)
	maps.Copy(minimumLoadoutNecessary, currentMinimumLoadout)
	for color, drawnNumber := range draw {
		if minimumLoadoutNecessary[color] < drawnNumber {
			minimumLoadoutNecessary[color] = drawnNumber
		}
	}

	return minimumLoadoutNecessary
}

func calculatePower(minimumLoadoutNecessary map[string]int) uint32 {
	var power uint32 = 1
	for _, number := range minimumLoadoutNecessary {
		power *= uint32(number)
	}
	return power
}

func solve(file *os.File) uint32 {
	scanner := bufio.NewScanner(file)

	var sum uint32 = 0
	for scanner.Scan() {
		line := scanner.Text()
		minimumLoadoutNecessary := map[string]int{
			"red":   0,
			"green": 0,
			"blue":  0,
		}
		for _, gameDraw := range strings.Split(line, ";") {
			parsedDraw := parseDraw(gameDraw)
			minimumLoadoutNecessary = determineMinimumLoadout(minimumLoadoutNecessary, parsedDraw)
		}
		sum += calculatePower(minimumLoadoutNecessary)
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
