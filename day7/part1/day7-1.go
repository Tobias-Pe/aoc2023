package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type CardValue int

func newCardValue(character rune) (CardValue, error) {
	switch character {
	case 'A':
		return A, nil
	case 'K':
		return K, nil
	case 'Q':
		return Q, nil
	case 'J':
		return J, nil
	case 'T':
		return T, nil
	case '9':
		return n9, nil
	case '8':
		return n8, nil
	case '7':
		return n7, nil
	case '6':
		return n6, nil
	case '5':
		return n5, nil
	case '4':
		return n4, nil
	case '3':
		return n3, nil
	case '2':
		return n2, nil
	default:
		return -1, fmt.Errorf("there is no CardValue for rune %c", character)
	}
}

const (
	n2 CardValue = iota
	n3
	n4
	n5
	n6
	n7
	n8
	n9
	T
	J
	Q
	K
	A
)

type HandType int

const (
	HighCard HandType = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

type Hand struct {
	cards []CardValue
	bid   int
}

func (hand Hand) getHandValue() HandType {
	var cardsCount = map[int]int{}
	for _, card := range hand.cards {
		valueOfCard := int(card)
		count, ok := cardsCount[valueOfCard]
		if !ok {
			count = 0
			cardsCount[valueOfCard] = count
		}
		cardsCount[valueOfCard] = count + 1
	}

	maxCount := 0
	secondMaxCount := 0
	for _, count := range cardsCount {
		if count == 5 {
			return FiveOfAKind
		} else if count == 4 {
			return FourOfAKind
		}

		if count >= maxCount {
			secondMaxCount = maxCount
			maxCount = count
		} else if count > secondMaxCount {
			secondMaxCount = count
		}
	}

	if secondMaxCount <= 1 {
		if maxCount == 3 {
			return ThreeOfAKind
		} else if maxCount == 2 {
			return OnePair
		}
		return HighCard
	}

	if maxCount == 3 {
		return FullHouse
	}

	return TwoPair
}

func main() {
	start := time.Now()

	lines := readFile("day7/part1/input.txt")
	hands, err := parseHands(lines)
	if err != nil {
		fmt.Println(err)
		return
	}

	handTypeToHands := groupHandsByHandType(hands)
	totalWinnings := calculateTotalWinnings(handTypeToHands)
	fmt.Println("[Part1] Total Winnings:", totalWinnings)

	fmt.Println("Finished in", time.Since(start))
}

func calculateTotalWinnings(handTypeToHands map[HandType][]Hand) int64 {
	rank := int64(1)
	totalWinnings := int64(0)
	for i := 0; i <= int(FiveOfAKind); i++ {
		handType := HandType(i)
		slices.SortFunc(handTypeToHands[handType], func(hand1, hand2 Hand) int {
			for handIndex := range hand1.cards {
				cardsCompareResult := cmp.Compare(int(hand1.cards[handIndex]), int(hand2.cards[handIndex]))
				if cardsCompareResult == 0 {
					continue
				}
				return cardsCompareResult
			}

			fmt.Println("Equal Hands ?")
			return 0
		})

		for _, hand := range handTypeToHands[handType] {
			totalWinnings += rank * int64(hand.bid)
			rank++
		}
	}
	return totalWinnings
}

func groupHandsByHandType(hands []Hand) map[HandType][]Hand {
	handTypeToHands := map[HandType][]Hand{}
	for _, hand := range hands {
		handValue := hand.getHandValue()
		hands, ok := handTypeToHands[handValue]
		if !ok {
			hands = []Hand{}
		}
		hands = append(hands, hand)
		handTypeToHands[handValue] = hands
	}
	return handTypeToHands
}

func readFile(file string) []string {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	split := strings.Split(lines, "\n")
	return split
}

func parseHands(lines []string) ([]Hand, error) {
	var hands []Hand
	for _, line := range lines {
		var hand = Hand{}
		handAndBid := strings.Fields(line)

		for _, cardSymbol := range handAndBid[0] {
			cardValue, err := newCardValue(cardSymbol)
			if err != nil {
				return nil, err
			}
			hand.cards = append(hand.cards, cardValue)
		}

		bid, err := strconv.Atoi(handAndBid[1])
		if err != nil {
			return nil, err
		}
		hand.bid = bid

		hands = append(hands, hand)
	}

	return hands, nil
}
