package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type pokerHand struct {
	cards string
	bid   int
}

func readFile(input *os.File) []pokerHand {
	scanner := bufio.NewScanner(input)
	var hands []pokerHand

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, " ")
		cards := splitLine[0]
		bid, _ := strconv.Atoi(splitLine[1])
		singleHand := pokerHand{cards: cards, bid: bid}
		hands = append(hands, singleHand)
	}

	return hands
}

func convertJokers(hand map[rune]int) map[rune]int {
	numberOfJokers := hand['J']
	hand['J'] = 0
	var cardWithMaximumNumber rune
	maximumNumber := -1
	for card, number := range hand {
		if number > maximumNumber {
			maximumNumber = number
			cardWithMaximumNumber = card
		}
	}
	hand[cardWithMaximumNumber] += numberOfJokers

	return hand
}

func determineHandScore(hand pokerHand) int {
	types := map[int][5]int{
		600: {5},
		500: {4},
		400: {3, 2},
		300: {3},
		200: {2, 2},
		100: {2},
	}
	orderedScores := [6]int{600, 500, 400, 300, 100}
	mappedHand := make(map[rune]int)
	for _, card := range hand.cards {
		mappedHand[card]++
	}
	mappedHand = convertJokers(mappedHand)
	reversedHand := make(map[int]rune)
	for card, number := range mappedHand {
		reversedHand[number] = card
	}
	var handScore int = 10
	for _, score := range orderedScores {
		conditions := types[score]
		foundType := true
		for _, condition := range conditions {
			if condition != 0 && reversedHand[condition] == 0 {
				foundType = false
				break
			}
		}
		if foundType {
			handScore = score
			break
		}
	}
	if handScore == 100 {
		// check two pairs manually
		numberOfTwos := 0
		for _, number := range mappedHand {
			if number == 2 {
				numberOfTwos++
			}
		}
		if numberOfTwos == 2 {
			handScore = 200
		}
	}
	return handScore
}

func compareHands(handA, handB pokerHand) bool {
	cards := map[rune]int{
		'A': 100,
		'K': 14,
		'Q': 13,
		'T': 10,
		'9': 9,
		'8': 8,
		'7': 7,
		'6': 6,
		'5': 5,
		'4': 4,
		'3': 3,
		'2': 2,
		'J': 1,
	}
	scoreA := determineHandScore(handA)
	scoreB := determineHandScore(handB)
	if scoreA == scoreB {
		for i := 0; i < len(handA.cards); i++ {
			cardA := rune(handA.cards[i])
			cardB := rune(handB.cards[i])
			if cardA != cardB {
				scoreA += cards[cardA]
				scoreB += cards[cardB]
				break
			}
		}
	}
	return scoreA < scoreB
}

func solve(hands []pokerHand) int {
	sort.Slice(hands[:], func(i, j int) bool {
		return compareHands(hands[i], hands[j])
	})
	totalWinnings := 0
	for ix, hand := range hands {
		totalWinnings += (ix + 1) * hand.bid
	}

	return totalWinnings
}

func main() {
	inputFile, _ := os.Open("input")

	hands := readFile(inputFile)
	inputFile.Close()

	// fmt.Println(determineHandScore(hand{cards: "AAAAA", bid: 100})) // 600
	// fmt.Println(determineHandScore(hand{cards: "KKKKA", bid: 100})) // 500
	// fmt.Println(determineHandScore(hand{cards: "KKKAA", bid: 100})) // 400
	// fmt.Println(determineHandScore(hand{cards: "AAAKQ", bid: 100})) // 300
	// fmt.Println(determineHandScore(hand{cards: "AAKKT", bid: 100})) // 200
	// fmt.Println(determineHandScore(hand{cards: "AAKQT", bid: 100})) // 100
	// fmt.Println(determineHandScore(hand{cards: "AKQT9", bid: 100})) // 10

	// fmt.Println(convertJokers(map[rune]int{'3': 2, 'T': 1, 'K': 1}))
	// fmt.Println(convertJokers(map[rune]int{'T': 1, 'J': 1, '5': 3}))
	// fmt.Println(convertJokers(map[rune]int{'K': 2, '6': 1, '7': 2}))
	// fmt.Println(convertJokers(map[rune]int{'K': 1, 'T': 2, 'J': 2}))
	// fmt.Println(convertJokers(map[rune]int{'Q': 3, 'J': 1, 'A': 1}))
	fmt.Println(solve(hands))
}
