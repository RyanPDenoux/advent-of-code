package solutions

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/ryanpdenoux/advent-of-code/utils"
)

var jokerVariant = flag.Bool(
	"joker-variant",
	false,
	"declare if using joker variation for Day7",
)

func Day7(file *os.File) {
	if !*jokerVariant {
		gameChoice()
	}

	parser := newCamelGameParser(file, *jokerVariant)
	game := parser.Parse()
	fmt.Printf("Game Winnings: %v\n", game.Winnings())
}

func gameChoice() {
	var choice string
	allowed := []string{"normal", "variant"}

	for !utils.Contains(allowed, choice) {
		fmt.Print("Pick a game type (normal, variant): ")
		fmt.Scan(&choice)
		if !utils.Contains(allowed, choice) {
			fmt.Print("  Please pick normal or variant\r")
		}
	}
	if choice == "variant" {
		*jokerVariant = true
	}
}

type CamelGame struct {
	length int
	head   *CamelHand
}

func (g CamelGame) String() string {
	var sb strings.Builder
	curr := g.head

	if curr == nil {
		return fmt.Sprint("No cards inserted for game")
	}

	for curr.Next != nil {
		sb.WriteString(curr.String())
		sb.WriteString("->")
		curr = curr.Next
	}

	sb.WriteString(curr.String())
	return sb.String()
}

func (g *CamelGame) InsertHand(hand *CamelHand) {
	g.length += 1
	curr := g.head

	if curr == nil {
		slog.Debug("First Node", "node", hand)
		g.head = hand
		return
	}

	if hand.Less(curr) {
		slog.Debug("Hand smaller than Head", "hand", hand, "head", curr)
		hand.Next = curr
		g.head = hand
		return
	}

	for curr.Next != nil {
		if hand.Less(curr.Next) {
			slog.Debug("Inserted Here", "point", fmt.Sprintf("...%v->%v->%v...", curr, hand, curr.Next))
			hand.Next = curr.Next
			curr.Next = hand
			return
		}
		curr = curr.Next
	}
	slog.Debug("Largest Hand", "game", fmt.Sprintf("...%v->%v", curr, hand))
	curr.Next = hand
	return
}

func (g *CamelGame) Winnings() int {
	var sum int
	var prev *CamelHand
	curr := g.head

	if curr == nil {
		return 0
	}

	i, j := 1, 1
	for curr.Next != nil {
		if prev != nil && curr.cards == prev.cards {
			sum += curr.bid * j
		} else {
			sum += curr.bid * i
			j = i
		}
		i++
		prev = curr
		curr = curr.Next
	}
	sum += curr.bid * g.length

	return sum
}

const (
	_ = iota
	HighCard
	OnePair
	TwoPairs
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

type CamelHand struct {
	cards    [5]CamelCard
	bid      int
	strength int
	Next     *CamelHand
}

func (h CamelHand) String() string {
	return fmt.Sprintf("%v(%d)|%d", h.cards, h.strength, h.bid)
}

func (h *CamelHand) Less(hand *CamelHand) bool {
	if h.strength < hand.strength {
		return true
	}
	if h.strength > hand.strength {
		return false
	}
	for i := 0; i < len(h.cards); i++ {
		ourCard, theirCard := h.cards[i], hand.cards[i]
		if ourCard.rank < theirCard.rank {
			return true
		}
		if ourCard.rank > theirCard.rank {
			return false
		}
	}
	return false
}

func newCamelHand(bs []byte, variant bool) *CamelHand {
	hand := &CamelHand{}
	hand.cards = [5]CamelCard{}

	for i, b := range bs {
		hand.cards[i] = newCamelCard(b, variant)
	}
	matchedCards := hand.matchCards(variant)
	hand.determineType(matchedCards)

	return hand
}

// REVIEW following two functions probably fit better on a "Rules" struct
func (h *CamelHand) determineType(matchedCards []int) {
	if utils.Contains(matchedCards, 5) {
		h.strength = FiveOfAKind
		return
	}

	if utils.Contains(matchedCards, 4) {
		h.strength = FourOfAKind
		return
	}

	if utils.Contains(matchedCards, 3) && utils.Contains(matchedCards, 2) {
		h.strength = FullHouse
		return
	}

	if utils.Contains(matchedCards, 3) {
		h.strength = ThreeOfAKind
		return
	}

	if utils.Contains(matchedCards, 2) && len(matchedCards) == 3 {
		h.strength = TwoPairs
		return
	}

	if utils.Contains(matchedCards, 2) {
		h.strength = OnePair
		return
	}

	h.strength = HighCard
	return
}

func (h *CamelHand) matchCards(joker bool) []int {
	var max int
	var maxCard CamelCard
	matchedCards := make(map[CamelCard]int)
	occurences := []int{}

	for i := 0; i < len(h.cards); i++ {
		matchedCards[h.cards[i]] += 1
	}

	// joker variant game
	if joker {
		for s, v := range matchedCards {
			if v > max && s != Joker{
				max = v
				maxCard = s
			}
		}
		val, ok := matchedCards[Joker]
		if ok {
			matchedCards[maxCard] += val
			delete(matchedCards, Joker)
		}
	}

	for _, count := range matchedCards {
		occurences = append(occurences, count)
	}

	return occurences
}

type CamelCard struct {
	rank   int
	symbol rune
}

func (c CamelCard) String() string {
	return string(c.symbol)
}

var (
	Ace   = CamelCard{14, 'A'}
	King  = CamelCard{13, 'K'}
	Queen = CamelCard{12, 'Q'}
	Jack  = CamelCard{11, 'J'}
	Ten   = CamelCard{10, 'T'}
	Nine  = CamelCard{9, '9'}
	Eight = CamelCard{8, '8'}
	Seven = CamelCard{7, '7'}
	Six   = CamelCard{6, '6'}
	Five  = CamelCard{5, '5'}
	Four  = CamelCard{4, '4'}
	Three = CamelCard{3, '3'}
	Two   = CamelCard{2, '2'}
	Joker = CamelCard{1, 'J'}
)

func newCamelCard(b byte, jokerRule bool) CamelCard {
	switch b {
	case 'A':
		return Ace
	case 'K':
		return King
	case 'Q':
		return Queen
	case 'J':
		if jokerRule {
			return Joker
		}
		return Jack
	case 'T':
		return Ten
	case '9':
		return Nine
	case '8':
		return Eight
	case '7':
		return Seven
	case '6':
		return Six
	case '5':
		return Five
	case '4':
		return Four
	case '3':
		return Three
	case '2':
		return Two
	}
	log.Fatalf("Not a valid CamelCard: %v\n", b)
	return CamelCard{}
}

type CamelGameParser struct {
	scanner *bufio.Scanner
	variant bool
}

func newCamelGameParser(file *os.File, variant bool) *CamelGameParser {
	p := &CamelGameParser{}
	p.scanner = bufio.NewScanner(file)
	p.variant = variant
	return p
}

func (p *CamelGameParser) Parse() *CamelGame {
	game := &CamelGame{}

	for p.scanner.Scan() {
		line := p.scanner.Text()
		hand := p.parseHand(line)
		game.InsertHand(hand)
	}

	slog.Debug("Parsed Game", "game", game)
	return game
}

func (p *CamelGameParser) parseHand(line string) *CamelHand {
	fields := strings.Fields(line)
	hand := newCamelHand([]byte(fields[0]), p.variant)
	bid, err := strconv.Atoi(fields[1])
	if err != nil {
		log.Fatal("No good")
	}
	hand.bid = bid
	return hand
}
