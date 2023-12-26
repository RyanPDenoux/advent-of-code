package solutions

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ryanpdenoux/advent-of-code/utils"
)

func Day4(file *os.File) {
	var points int
	var count int

	cards := make(CopyMap)
	scanner := bufio.NewScanner(file)
	scoring := &ValueScoring{scoringBase: 2}
	parser := newGameParser(":", "|")
	for i := 1; scanner.Scan(); i++ {
		cards[i] += 1
		line := scanner.Text()
		game := parser.ParseGame(line)
		score := game.scoreGame(scoring)
		cards.insertCopies(game, i)
		points += score
		count = i
	}
	fmt.Printf("Value of Scratchcards: %v\n", points)
	fmt.Printf("Count of all cards: %v\t%v\n", cards.sumValues(count), cards)
}

type Scorable interface {
	ScoreGame(Set) int
}

type CopyMap map[int]int

func (m CopyMap) insertCopies(game *ScratchGame, offset int) {
	numMatches :=  len(game.matches)

	for _, i := range utils.MakeRange(1, numMatches+1) {
		pos := offset+i
		scaler := m[offset]
		m[pos] += scaler
	}
	slog.Debug("Updated map", "map", m)
}

func (m CopyMap) sumValues(end int) int {
	var sum int

	for i := 1; i <= end; i++ {
		sum += m[i]
	}

	return sum
}

type Set map[int]bool

func newSetFromStrSlice(slice []string) Set {
	set := Set{}

	for _, element := range slice {
		if element == "" {
			continue
		}
		i, err := strconv.Atoi(element)
		if err != nil {
			log.Fatalf("Could not convert item %v: %v\n", element, err)
		}
		set[i] = true
	}

	slog.Debug("Created Integer Slice", "slice", slice)
	return set
}

func (this Set) Union(other Set) Set {
	union := Set{}

	for k := range this {
		union[k] = true
	}

	for k := range other {
		union[k] = true
	}

	return union
}

func (this Set) Intersect(other Set) Set {
	intersection := Set{}

	for k := range this {
		if other[k] {
			intersection[k] = true
		}
	}

	return intersection
}

func (this Set) String() string {
	elements := []int{}
	for k := range this {
		elements = append(elements, k)
	}
	return fmt.Sprintf("Set%v", elements)
}

type GameParser struct {
	raw         string
	headerDelim string
	gameDelim   string
}

func newGameParser(headerChar, gameChar string) *GameParser {
	s := &GameParser{
		headerDelim: headerChar,
		gameDelim:   gameChar,
	}
	return s
}

func (s *GameParser) ParseGame(raw string) *ScratchGame {
	game := &ScratchGame{raw: raw}

	_, gameData := s.parseRaw(raw)
	game.winningNums, game.playerNums = s.parseGame(gameData)
	game.matches = game.winningNums.Intersect(game.playerNums)
	slog.Debug("Parsed game", "game", game)
	return game
}

func (s *GameParser) parseRaw(raw string) (string, string) {
	parsed := strings.Split(raw, s.headerDelim)
	if len(parsed) != 2 {
		log.Fatalf("Raw game input is incompatible: %v\n", raw)
	}

	return parsed[0], parsed[1]
}

func (s *GameParser) parseGame(game string) (Set, Set) {
	parsed := strings.Split(game, s.gameDelim)
	if len(parsed) != 2 {
		log.Fatalf("Game not properly formatted: %v\n", game)
	}

	winners := newSetFromStrSlice(s.splitGame(parsed[0]))
	players := newSetFromStrSlice(s.splitGame(parsed[1]))

	return winners, players
}

func (s *GameParser) splitGame(game string) []string {
	return strings.Fields(game)
}

type ScratchGame struct {
	raw         string
	winningNums Set
	playerNums  Set
	matches     Set
}

func (s *ScratchGame) scoreGame(scoringSystem *ValueScoring) int {
	if len(s.matches) == 0 {
		return 0
	}

	score := scoringSystem.ScoreGame(s.matches)
	slog.Debug("Calculated score for game","matches", s.matches, "score", score, "game", s)
	return score
}

type ValueScoring struct {
	scoringBase int
}

func (s *ValueScoring) ScoreGame(matches Set) int {
	numMatches := len(matches)
	score := int(math.Pow(float64(s.scoringBase), float64(numMatches - 1)))
	return score
}
